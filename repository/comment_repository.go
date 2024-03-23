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

	tx, err := conn.Begin(ctx)
	if err != nil {
		return customErr.ErrorInternalServer
	}

	SQL_GET_POST_ID := "SELECT EXISTS (SELECT 1 FROM posts WHERE post_id = $1)"
	_, err = tx.Exec(ctx, SQL_GET_POST_ID, comment.PostId)
	if err != nil {
		tx.Rollback(ctx)
		if err == pgx.ErrNoRows {
			return customErr.ErrorNotFound
		} else {
			return customErr.ErrorInternalServer
		}
	}

	SQL_GET_FRIEND := "SELECT EXISTS (SELECT 1 FROM users u " +
		"JOIN friends f ON u.user_id = f.user_id_requester " +
		"WHERE f.user_id_accepter = $1)"
	var isFriend bool

	err = tx.QueryRow(ctx, SQL_GET_FRIEND, comment.UserId).Scan(&isFriend)
	if err != nil {
		tx.Rollback(ctx)
		if err == pgx.ErrNoRows {
			return customErr.ErrorBadRequest
		} else {
			return customErr.ErrorInternalServer
		}
	}

	if !isFriend {
		tx.Rollback(ctx)
		return customErr.ErrorBadRequest
	}

	SQL_INSERT := "INSERT INTO comments(post_id, comment, user_id) values ($1, $2, $3)"

	result, err := tx.Exec(ctx, SQL_INSERT, comment.PostId, comment.Comment, comment.UserId)
	if err != nil {
		tx.Rollback(ctx)
		if err == pgx.ErrNoRows {
			return customErr.ErrorBadRequest
		} else {
			return customErr.ErrorInternalServer
		}
	}

	if result.RowsAffected() == 0 {
		tx.Rollback(ctx)
		return customErr.ErrorBadRequest
	}

	tx.Commit(ctx)

	return nil
}
