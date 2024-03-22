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

	SQL_GET_POST_ID := "SELECT EXISTS (SELECT 1 FROM posts WHERE post_id = $1)"
	var isPostExists bool

	err = conn.QueryRow(ctx, SQL_GET_POST_ID, comment.PostId).Scan(&isPostExists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return customErr.ErrorBadRequest
		} else {
			return customErr.ErrorInternalServer
		}
	}

	if !isPostExists {
		return customErr.ErrorNotFound
	}

	SQL_INSERT := "INSERT INTO comments(post_id, comment, user_id) " +
		"SELECT $1, $2, $3 " +
		"FROM posts p " +
		"JOIN users u ON p.user_id = u.user_id " +
		"JOIN friends f ON u.user_id = f.user_id_requester " +
		"WHERE f.user_id_accepter = $3 " +
		"RETURNING comment_id"

	result, err := conn.Exec(ctx, SQL_INSERT, comment.PostId, comment.Comment, comment.UserId)
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
