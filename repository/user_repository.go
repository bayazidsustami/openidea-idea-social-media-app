package repository

import (
	"context"
	"database/sql"
	"openidea-idea-social-media-app/customErr"
	user_model "openidea-idea-social-media-app/models/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Register(ctx context.Context, tx pgx.Tx, user user_model.User) (user_model.User, error)
	Login(ctx context.Context, conn *pgxpool.Pool, user user_model.User) (user_model.User, error)
	UpdateEmail(ctx context.Context, conn *pgxpool.Pool, userId string, email string) error
	UpdatePhone(ctx context.Context, conn *pgxpool.Pool, userId string, phone string) error
	Update(ctx context.Context, conn *pgxpool.Pool, user user_model.User) error
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
		SQL_INSERT = "INSERT INTO users(user_id, email, password, name) values (gen_random_uuid(), $1, $2, $3) " +
			"ON CONFLICT(email) " +
			"DO NOTHING RETURNING user_id"
		emailOrPhone = user.Email.String
	} else {
		SQL_INSERT = "INSERT INTO users(user_id, phone, password, name) values (gen_random_uuid(), $1, $2, $3) " +
			"ON CONFLICT(phone) " +
			"DO NOTHING RETURNING user_id"
		emailOrPhone = user.Phone.String
	}

	var idUser string
	err := tx.QueryRow(ctx, SQL_INSERT, emailOrPhone, user.Password, user.Name).Scan(&idUser)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user_model.User{}, customErr.ErrorConflict
		} else {
			return user_model.User{}, customErr.ErrorInternalServer
		}

	}

	user.UserId = idUser
	tx.Commit(ctx)
	return user, nil
}

func (repository *UserRepositoryImpl) Login(ctx context.Context, conn *pgxpool.Pool, user user_model.User) (user_model.User, error) {
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

func (repository *UserRepositoryImpl) UpdateEmail(ctx context.Context, conn *pgxpool.Pool, userId string, email string) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return customErr.ErrorInternalServer
	}
	defer tx.Rollback(ctx)

	GET_EMAIL := "SELECT email FROM users u WHERE u.user_id = $1"

	var resEmail sql.NullString
	err = tx.QueryRow(ctx, GET_EMAIL, userId).Scan(&resEmail)
	if err != nil {
		if err == pgx.ErrNoRows {
			return customErr.ErrorBadRequest
		} else {
			return customErr.ErrorInternalServer
		}
	}

	if resEmail.String != "" {
		return customErr.ErrorBadRequest
	}

	GET_IS_EXISTED := "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)"
	var isExistsEmail bool
	err = tx.QueryRow(ctx, GET_IS_EXISTED, email).Scan(&isExistsEmail)
	if err != nil {
		if err == pgx.ErrNoRows {
			return customErr.ErrorBadRequest
		} else {
			return customErr.ErrorInternalServer
		}
	}

	if isExistsEmail {
		return customErr.ErrorConflict
	}

	UPDATE_EMAIL := `
		UPDATE users
		SET email = $1
		WHERE user_id = $2
		AND (email IS NULL OR email = '')
	`
	res, err := tx.Exec(ctx, UPDATE_EMAIL, email, userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return customErr.ErrorBadRequest
		} else {
			return customErr.ErrorInternalServer
		}
	}

	if res.RowsAffected() == 0 {
		return customErr.ErrorNotFound
	}

	tx.Commit(ctx)

	return nil
}

func (repository *UserRepositoryImpl) UpdatePhone(ctx context.Context, conn *pgxpool.Pool, userId string, phone string) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return customErr.ErrorInternalServer
	}
	defer tx.Rollback(ctx)

	GET_PHONE := "SELECT phone FROM users u WHERE u.user_id = $1"

	var resPhone sql.NullString
	err = tx.QueryRow(ctx, GET_PHONE, userId).Scan(&resPhone)
	if err != nil {
		if err == pgx.ErrNoRows {
			return customErr.ErrorBadRequest
		} else {
			return customErr.ErrorInternalServer
		}
	}

	if resPhone.String != "" {
		return customErr.ErrorBadRequest
	}

	GET_IS_EXISTED := "SELECT EXISTS (SELECT 1 FROM users WHERE phone = $1)"
	var isExistsPhone bool
	err = tx.QueryRow(ctx, GET_IS_EXISTED, phone).Scan(&isExistsPhone)
	if err != nil {
		if err == pgx.ErrNoRows {
			return customErr.ErrorBadRequest
		} else {
			return customErr.ErrorInternalServer
		}
	}

	if isExistsPhone {
		return customErr.ErrorConflict
	}

	UPDATE_EMAIL := `
		UPDATE users
		SET phone = $1
		WHERE user_id = $2
		AND (phone IS NULL OR phone = '')
	`
	res, err := tx.Exec(ctx, UPDATE_EMAIL, phone, userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return customErr.ErrorBadRequest
		} else {
			return customErr.ErrorInternalServer
		}
	}

	if res.RowsAffected() == 0 {
		return customErr.ErrorNotFound
	}

	tx.Commit(ctx)

	return nil
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, conn *pgxpool.Pool, user user_model.User) error {
	UPDATE_ACC := "UPDATE users " +
		"SET name = $1, image_url = $2, updated_at = CURRENT_TIMESTAMP " +
		"WHERE user_id = $3"

	res, err := conn.Exec(ctx, UPDATE_ACC, user.Name, user.ImageUrl, user.UserId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return customErr.ErrorBadRequest
		} else {
			return customErr.ErrorInternalServer
		}
	}

	if res.RowsAffected() == 0 {
		return customErr.ErrorNotFound
	}

	return nil
}
