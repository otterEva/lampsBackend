package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/otterEva/lamps/users_service/logs"
	"github.com/otterEva/lamps/users_service/settings"
)

func GenerateToken(admin bool, userId uint) (string, error) {
	secret := []byte(settings.Config.SECRET)
	method := jwt.SigningMethodHS256

	claims := jwt.MapClaims{
		"admin":  admin,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 168).Unix(),
	}

	token, err := jwt.NewWithClaims(method, claims).SignedString(secret)
	if err != nil {
		logs.Logger.Error("error while creating token", "data", claims)
		return "", err
	}

	logs.Logger.Debug("token has been created")
	return token, nil
}
