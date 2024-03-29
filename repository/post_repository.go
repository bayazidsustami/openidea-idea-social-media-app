package repository

import (
	"context"
	"openidea-idea-social-media-app/customErr"
	post_model "openidea-idea-social-media-app/models/post"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository interface {
	Create(ctx context.Context, post post_model.Post, userId string) error
	GetAll(ctx context.Context, filters post_model.PostFilters) ([]post_model.Post, int, error)
}

type PostRepositoryImpl struct {
	DBPool *pgxpool.Pool
}

func NewPostRepository(DBPool *pgxpool.Pool) PostRepository {
	return &PostRepositoryImpl{
		DBPool: DBPool,
	}
}

func (repository *PostRepositoryImpl) Create(ctx context.Context, post post_model.Post, userId string) error {

	SQL_INSERT := "INSERT INTO posts(post_html, tags, user_id) values ($1, $2, $3)"

	res, err := repository.DBPool.Exec(ctx, SQL_INSERT, post.PostHtml, post.Tags, userId)
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

func (repository *PostRepositoryImpl) GetAll(ctx context.Context, filters post_model.PostFilters) ([]post_model.Post, int, error) {
	query := filters.BuildQuery()

	rows, err := repository.DBPool.Query(ctx, query)
	if err != nil {
		return nil, 0, customErr.ErrorInternalServer
	}
	defer rows.Close()

	var posts []post_model.Post
	var totalPosts int

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
			&totalPosts,
		)

		if err != nil {
			return nil, 0, customErr.ErrorInternalServer
		}

		posts = append(posts, post)
	}

	return posts, totalPosts, nil
}
