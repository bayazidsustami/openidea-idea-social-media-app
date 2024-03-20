package repository

import (
	"context"
	"openidea-idea-social-media-app/customErr"
	comment_model "openidea-idea-social-media-app/models/comment"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepository interface {
	Create(ctx context.Context, request comment_model.Comment) error
}

type CommentRepositoryImpl struct {
	DBPool *pgxpool.Pool
}

func NewCommentRepository(DBPool *pgxpool.Pool) CommentRepository {
	return &CommentRepositoryImpl{
		DBPool: DBPool,
	}
}

func (repository *CommentRepositoryImpl) Create(ctx context.Context, comment comment_model.Comment) error {
	conn, err := repository.DBPool.Acquire(ctx)
	if err != nil {
		return customErr.ErrorInternalServer
	}
	defer conn.Release()

	// TODO: Need to add validation to only comment on user's friend's post
	SQL_INSERT := "INSERT INTO comments(post_id, user_id, comment) values ($1, $2, $3) RETURNING comment_id"

	result, err := conn.Exec(ctx, SQL_INSERT, comment.PostId, comment.UserId, comment.Comment)
	if err != nil {
		if err == pgx.ErrNoRows {
			return customErr.ErrorBadRequest
		} else {
			return customErr.ErrorInternalServer
		}
	}

	if result.RowsAffected() == 0 {
		return customErr.ErrorBadRequest
	}

	return nil
}
