package service

import (
	"context"
	"openidea-idea-social-media-app/customErr"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var tokenExpDuration = time.Now().Add(time.Minute * 30).Unix()

const JWT_TOKEN_LOCALS = "userInfo"

type AuthService interface {
	GenerateToken(ctx context.Context, userId int) (string, error)
	GetValidUser(ctx *fiber.Ctx) (int, error)
}

type AuthServiceImpl struct {
}

func NewAuthService() AuthService {
	return &AuthServiceImpl{}
}

func (service *AuthServiceImpl) GenerateToken(ctx context.Context, userId int) (string, error) {

	claims := jwt.MapClaims{
		"username": userId,
		"exp":      tokenExpDuration,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		return "", customErr.ErrorInternalServer
	}

	return signedToken, nil
}

func (service *AuthServiceImpl) GetValidUser(ctx *fiber.Ctx) (int, error) {
	userInfo := ctx.Locals(JWT_TOKEN_LOCALS).(*jwt.Token)
	claims := userInfo.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)

	return int(userId), nil
}
