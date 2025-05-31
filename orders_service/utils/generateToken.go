package utils

import (
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt"
	"github.com/otterEva/lamps/orders_service/settings"
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
		log.Debug(err.Error())
		return "", err
	}

	return token, nil
}
