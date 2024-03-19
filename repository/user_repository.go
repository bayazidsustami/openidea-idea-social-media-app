package repository

import (
	"context"
	"openidea-idea-social-media-app/customErr"
	user_model "openidea-idea-social-media-app/models/user"

	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Register(ctx context.Context, tx pgx.Tx, user user_model.User) (user_model.User, error)
}

type UserRepositoryImpl struct {
}

func New() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Register(ctx context.Context, tx pgx.Tx, user user_model.User) (user_model.User, error) {
	var SQL_INSERT string
	var emailOrPhone string
	if user.Email != "" {
		SQL_INSERT = "INSERT INTO users(email, password, name) values ($1, $2, $3) " +
			"ON CONFLICT(email) " +
			"DO NOTHING RETURNING user_id"
		emailOrPhone = user.Email
	} else {
		SQL_INSERT = "INSERT INTO users(phone, password, name) values ($1, $2, $3) " +
			"ON CONFLICT(phone) " +
			"DO NOTHING RETURNING user_id"
		emailOrPhone = user.Phone
	}

	var idUser int
	err := tx.QueryRow(ctx, SQL_INSERT, emailOrPhone, user.Password, user.Name).Scan(&idUser)
	if err != nil {
		return user_model.User{}, customErr.ErrorConflict
	}

	user.UserId = idUser
	tx.Commit(ctx)
	return user, nil
}
