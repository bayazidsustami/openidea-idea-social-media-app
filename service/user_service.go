package service

import (
	"context"
	"log"
	"openidea-idea-social-media-app/customErr"
	user_model "openidea-idea-social-media-app/models/user"
	"openidea-idea-social-media-app/repository"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	Register(ctx context.Context, request user_model.UserRegisterRequest) user_model.UserRegisterResponse[user_model.UserData]
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	Validator      *validator.Validate
	DBPool         *pgxpool.Pool
}

func New(
	userRepository repository.UserRepository,
	validator *validator.Validate,
	dbPool *pgxpool.Pool,
) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		Validator:      validator,
		DBPool:         dbPool,
	}
}

func (service *UserServiceImpl) Register(ctx context.Context, request user_model.UserRegisterRequest) user_model.UserRegisterResponse[user_model.UserData] {
	err := service.Validator.Struct(request)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := service.DBPool.Acquire(ctx)
	if err != nil {
		log.Fatal(customErr.ErrorInternalServer)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatal(customErr.ErrorInternalServer)
	}
	defer tx.Rollback(ctx)

	user := user_model.User{}

	userResult := service.UserRepository.Register(ctx, tx, user)

	if request.CredentialType == "email" {
		return user_model.UserRegisterResponse[user_model.UserData]{
			Message: "User registered successfully",
			Data: &user_model.UserEmailDataResponse{
				Email: userResult.Email,
				Name:  userResult.Name,
			},
		}
	} else {
		return user_model.UserRegisterResponse[user_model.UserData]{
			Message: "User registered successfully",
			Data: &user_model.UserPhoneDataResponse{
				Phone: user.Phone,
				Name:  userResult.Name,
			},
		}
	}

}
