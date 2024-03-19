package service

import (
	"context"
	"log"
	"openidea-idea-social-media-app/customErr"
	user_model "openidea-idea-social-media-app/models/user"
	"openidea-idea-social-media-app/repository"
	"openidea-idea-social-media-app/security"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	Register(ctx context.Context, request user_model.UserRegisterRequest) (user_model.UserRegisterResponse[user_model.UserData], error)
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	Validator      *validator.Validate
	DBPool         *pgxpool.Pool
	AuthService    AuthService
}

func NewUserService(
	userRepository repository.UserRepository,
	validator *validator.Validate,
	dbPool *pgxpool.Pool,
	authService AuthService,
) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		Validator:      validator,
		DBPool:         dbPool,
	}
}

func (service *UserServiceImpl) Register(ctx context.Context, request user_model.UserRegisterRequest) (user_model.UserRegisterResponse[user_model.UserData], error) {
	err := service.Validator.Struct(request)
	if err != nil {
		log.Println(err)
		return user_model.UserRegisterResponse[user_model.UserData]{}, customErr.ErrorBadRequest
	}

	conn, err := service.DBPool.Acquire(ctx)
	if err != nil {
		return user_model.UserRegisterResponse[user_model.UserData]{}, customErr.ErrorInternalServer
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return user_model.UserRegisterResponse[user_model.UserData]{}, customErr.ErrorInternalServer
	}
	defer tx.Rollback(ctx)

	hashedPass, err := security.GenerateHashedPassword(request.Password)

	if err != nil {
		return user_model.UserRegisterResponse[user_model.UserData]{}, err
	}

	user := user_model.User{
		Password: hashedPass,
	}

	if request.CredentialType == "email" {
		user.Email = request.CredentialValue
	} else {
		user.Password = request.CredentialValue
	}

	userResult, err := service.UserRepository.Register(ctx, tx, user)
	if err != nil {
		return user_model.UserRegisterResponse[user_model.UserData]{}, err
	}

	token, err := service.AuthService.GenerateToken(ctx, userResult.UserId)
	if err != nil {
		tx.Rollback(ctx)
		return user_model.UserRegisterResponse[user_model.UserData]{}, err
	}

	if request.CredentialType == "email" {
		return user_model.UserRegisterResponse[user_model.UserData]{
			Message: "User registered successfully",
			Data: &user_model.UserEmailDataResponse{
				Email:       userResult.Email,
				Name:        userResult.Name,
				AccessToken: token,
			},
		}, nil
	} else {
		return user_model.UserRegisterResponse[user_model.UserData]{
			Message: "User registered successfully",
			Data: &user_model.UserPhoneDataResponse{
				Phone:       userResult.Phone,
				Name:        userResult.Name,
				AccessToken: token,
			},
		}, nil
	}

}
