package repository

import (
	"context"
	"openidea-idea-social-media-app/customErr"
	user_model "openidea-idea-social-media-app/models/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Register(ctx context.Context, tx pgx.Tx, user user_model.User) (user_model.User, error)
	Login(ctx context.Context, conn *pgxpool.Conn, user user_model.User) (user_model.User, error)
}

type UserRepositoryImpl struct {
}

func New() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Register(ctx context.Context, tx pgx.Tx, user user_model.User) (user_model.User, error) {
	var SQL_INSERT string
	var emailOrPhone string
	if user.Email.String != "" {
		SQL_INSERT = "INSERT INTO users(email, password, name) values ($1, $2, $3) " +
			"ON CONFLICT(email) " +
			"DO NOTHING RETURNING user_id"
		emailOrPhone = user.Email.String
	} else {
		SQL_INSERT = "INSERT INTO users(phone, password, name) values ($1, $2, $3) " +
			"ON CONFLICT(phone) " +
			"DO NOTHING RETURNING user_id"
		emailOrPhone = user.Phone.String
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

func (repository *UserRepositoryImpl) Login(ctx context.Context, conn *pgxpool.Conn, user user_model.User) (user_model.User, error) {
	var SQL_GET_USER string
	var emailOrPhone string
	if user.Email.String != "" {
		SQL_GET_USER = "SELECT user_id, email, phone, name, password FROM users WHERE email=$1"
		emailOrPhone = user.Email.String
	} else {
		SQL_GET_USER = "SELECT user_id, email, phone, name, password FROM users WHERE phone=$1"
		emailOrPhone = user.Phone.String
	}

	result := user_model.User{}
	err := conn.QueryRow(ctx, SQL_GET_USER, emailOrPhone).Scan(
		&result.UserId,
		&result.Email,
		&result.Phone,
		&result.Name,
		&result.Password,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return user_model.User{}, customErr.ErrorNotFound
		} else {
			return user_model.User{}, customErr.ErrorInternalServer
		}
	}

	return result, nil
}
