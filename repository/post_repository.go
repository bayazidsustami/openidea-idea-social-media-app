package repository

import (
	"context"
	post_model "openidea-idea-social-media-app/models/post"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository interface {
	Create(ctx context.Context, post post_model.Post) error
	GetAll(ctx context.Context) ([]post_model.Post, error)
}

type PostRepositoryImpl struct {
	DBPool *pgxpool.Pool
}

func New(DBPool *pgxpool.Pool) PostRepository {
	return &PostRepositoryImpl{
		DBPool: DBPool,
	}
}

func (repository *PostRepositoryImpl) Create(ctx context.Context, post post_model.Post) error {
	conn, err := repository.DBPool.Acquire(ctx)
	if err != nil {
		return err
	}

	SQL_INSERT := "INSERT INTO posts(post_html, tags, user_id) values ($1, $2, $3) RETURNING post_id"

	_, err = conn.Exec(ctx, SQL_INSERT, post.PostHtml, post.Tags)
	if err != nil {
		return err
	}

	return nil
}

func (repository *PostRepositoryImpl) GetAll(ctx context.Context) ([]post_model.Post, error) {
	conn, err := repository.DBPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	SQL_GET := "SELECT p.post_id, p.post_html, p.tags, p.created_at, " +
		"u.user_id, u.name, u.image_url, u.created_at, " +
		"jsonb_agg(jsonb_build_object(" +
		"'comment', c.comment," +
		"'createdAt', c.created_at" +
		// TODO: Adjust this
		// "'creator', SELECT jsonb_build_object('userId', cu.user_id, 'name', cu.name, 'imageUrl', cu.image_url, 'createdAt', cu.created_at) AS 'creator'" +
		")) AS 'comments' " +
		"FROM posts p " +
		"LEFT JOIN users u ON p.user_id = u.user_id " +
		"LEFT JOIN comments c ON p.post_id = c.post_id " +
		"GROUP BY p.post_id"

	rows, err := conn.Query(ctx, SQL_GET)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fiber.NewError(fiber.StatusNotFound, "Not Found")
	}

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
			&post.Comments,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
