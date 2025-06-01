package middlewares

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/otterEva/lamps/image_service/settings"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// блок работы с куки

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

		if !ok {
			c.ClearCookie("jwt")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "invalid token",
			})
		}

		// блок проверки прав

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()

			resp, err := http.Get(fmt.Sprintf("http://users_service:8001/%s/%s", userId, admin))
			if err != nil {
				fmt.Println("Ошибка запроса:", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				c.SendStatus(fiber.StatusForbidden)
			} else {
				c.Locals("userId", userId)
				c.Locals("admin", admin)
			}
		}()

		wg.Wait()

		return c.Next()
	}
}
