package middlewares

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx"
	"github.com/otterEva/lamps/image_service/settings"
)

func AuthMiddleware(ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {


		cookieToken := c.Cookies("jwt")

		if cookieToken == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		token, err := jwt.Parse(cookieToken, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(settings.Config.SECRET), nil
		})

		if err != nil || !token.Valid {
			c.ClearCookie("jwt")

			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "Unauthorized",
			})
		}

		claims := token.Claims.(jwt.MapClaims)

		userId, ok := claims["userId"]
		admin, ok := claims["admin"]
		
		log.Debug(userId, admin)

		if !ok {
			c.ClearCookie("jwt")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "invalid token",
			})
		}
		
		log.Debug(userId, admin)

		sql, args, err := sq.
			Select("1").
			From("Users").
			Where(sq.Eq{"id": userId, "admin": admin}).
			Limit(1).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			log.Debug(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var exists int

		err = dbClient.QueryRow(ctx, sql, args...).Scan(&exists)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.ClearCookie("jwt")
				return c.SendStatus(fiber.StatusUnauthorized)
			}
			log.Fatal(err)
		}

		c.Locals("userId", userId)
		c.Locals("admin", admin)

		return c.Next()
	}
}
