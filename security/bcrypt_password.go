package security

import (
	"openidea-idea-social-media-app/customErr"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func GenrateHashedPassword(reqPassword string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(reqPassword), viper.GetInt("BCRYPT_SALT"))
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}

func ComparePassword(hashedPassword string, reqPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(reqPassword))
	if err != nil {
		return customErr.ErrorBadRequest
	}

	return nil
}
