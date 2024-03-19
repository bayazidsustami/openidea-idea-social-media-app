package service

import (
	"context"
	"openidea-idea-social-media-app/customErr"
	user_model "openidea-idea-social-media-app/models/user"
	"openidea-idea-social-media-app/repository"
	"openidea-idea-social-media-app/security"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	Register(ctx context.Context, request user_model.UserRegisterRequest) (user_model.UserRegisterResponse[user_model.UserData], error)
	Login(ctx context.Context, request user_model.UserLoginRequest) (user_model.UserLoginResponse, error)
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
		AuthService:    authService,
	}
}

func (service *UserServiceImpl) Register(ctx context.Context, request user_model.UserRegisterRequest) (user_model.UserRegisterResponse[user_model.UserData], error) {
	err := service.Validator.Struct(request)
	if err != nil {
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

	var user user_model.User
	if request.CredentialType == "email" {
		user = user_model.User{
			Password: hashedPass,
			Name:     request.Name,
			Email:    request.CredentialValue,
		}
	} else {
		user = user_model.User{
			Password: hashedPass,
			Name:     request.Name,
			Phone:    request.CredentialValue,
		}
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

func (service *UserServiceImpl) Login(ctx context.Context, request user_model.UserLoginRequest) (user_model.UserLoginResponse, error) {
	err := service.Validator.Struct(request)
	if err != nil {
		return user_model.UserLoginResponse{}, customErr.ErrorBadRequest
	}

	conn, err := service.DBPool.Acquire(ctx)
	if err != nil {
		return user_model.UserLoginResponse{}, customErr.ErrorInternalServer
	}
	defer conn.Release()

	var user user_model.User
	if request.CredentialType == "email" {
		user = user_model.User{
			Email: request.CredentialValue,
		}
	} else {
		user = user_model.User{
			Phone: request.CredentialValue,
		}
	}

	result, err := service.UserRepository.Login(ctx, conn, user)
	if err != nil {
		return user_model.UserLoginResponse{}, err
	}

	err = security.ComparePassword(result.Password, request.Password)
	if err != nil {
		return user_model.UserLoginResponse{}, err
	}

	token, err := service.AuthService.GenerateToken(ctx, result.UserId)
	if err != nil {
		return user_model.UserLoginResponse{}, err
	}

	userResponse := user_model.UserLoginResponse{
		Message: "User logged successfully",
		Data: user_model.UserEmailPhoneDataResponse{
			Name:        result.Name,
			AccessToken: token,
			Email:       result.Email,
			Phone:       result.Phone,
		},
	}

	return userResponse, nil
}
