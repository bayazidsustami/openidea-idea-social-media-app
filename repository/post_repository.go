package repository

import (
	"context"
	"openidea-idea-social-media-app/customErr"
	post_model "openidea-idea-social-media-app/models/post"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository interface {
	Create(ctx context.Context, post post_model.Post, userId int) error
	GetAll(ctx context.Context, filters post_model.PostFilters) ([]post_model.Post, error)
}

type PostRepositoryImpl struct {
	DBPool *pgxpool.Pool
}

func NewPostRepository(DBPool *pgxpool.Pool) PostRepository {
	return &PostRepositoryImpl{
		DBPool: DBPool,
	}
}

func (repository *PostRepositoryImpl) Create(ctx context.Context, post post_model.Post, userId int) error {
	conn, err := repository.DBPool.Acquire(ctx)
	if err != nil {
		return customErr.ErrorInternalServer
	}

	SQL_INSERT := "INSERT INTO posts(post_html, tags, user_id) values ($1, $2, $3)"

	res, err := conn.Exec(ctx, SQL_INSERT, post.PostHtml, post.Tags, userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return customErr.ErrorBadRequest
		} else {
			return customErr.ErrorInternalServer
		}
	}

	if res.RowsAffected() == 0 {
		return customErr.ErrorBadRequest
	}

	return nil
}

// TODO: Populate filters to query
func (repository *PostRepositoryImpl) GetAll(ctx context.Context, filters post_model.PostFilters) ([]post_model.Post, error) {
	conn, err := repository.DBPool.Acquire(ctx)
	if err != nil {
		return nil, customErr.ErrorInternalServer
	}
	defer conn.Release()

	SQL_GET := "SELECT p.post_id, p.post_html, p.tags, p.created_at, " +
		"u.user_id, u.name, u.image_url, u.created_at, " +
		"(SELECT COUNT(*) FROM friends f WHERE p.user_id = f.user_id_requester), " +
		"jsonb_agg(jsonb_build_object(" +
		"'comment', c.comment," +
		"'createdAt', c.created_at," +
		`'creator', jsonb_build_object('userId', cu.user_id, 'name', cu.name, 'imageUrl', cu.image_url, 'friendCount', (SELECT COUNT(*) FROM friends cf WHERE cu.user_id = cf.user_id_requester), 'createdAt', to_char(cu.created_at, 'YYYY-MM-DD"T"HH24:MI:SSOF'))` +
		")) AS comments " +
		"FROM posts p " +
		"JOIN users u ON p.user_id = u.user_id " +
		"LEFT JOIN comments c ON p.post_id = c.post_id " +
		"JOIN users cu ON c.user_id = cu.user_id " +
		"GROUP BY p.post_id, u.user_id"

	rows, err := conn.Query(ctx, SQL_GET)
	if err != nil {
		return nil, customErr.ErrorInternalServer
	}
	defer rows.Close()

	var posts []post_model.Post
	for rows.Next() {
		post := post_model.Post{}

		err = rows.Scan(
			&post.PostId,
			&post.PostHtml,
			&post.Tags,
			&post.CreatedAt,
			&post.Creator.UserId,
			&post.Creator.Name,
			&post.Creator.ImageUrl,
			&post.Creator.CreatedAt,
			&post.Creator.FriendCount,
			&post.Comments,
		)

		if err != nil {
			return nil, customErr.ErrorInternalServer
		}

		posts = append(posts, post)
	}

	if len(posts) == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Not Found")
	}

	return posts, nil
}
