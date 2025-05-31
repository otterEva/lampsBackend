package handlers

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/otterEva/lamps/auth_service/schemas"
	"github.com/otterEva/lamps/auth_service/utils"
	"golang.org/x/crypto/bcrypt"
)

func AuthHandlers(route fiber.Router, dbClient *pgxpool.Pool, ctx context.Context) {
	// @Summary      Register user
	// @Description  register by email and password
	// @Tags         auth
	// @Accept       form-data
	// @Param        email     body  string  true  "Email"
	// @Param        password  body  string  true  "Password"
	// @Router       /auth/register [post]
	route.Post("/register", func(c *fiber.Ctx) error {
		user := &schemas.User{
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
		}

		if user.Email == "" || user.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email and password required",
			})
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Debug(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		sql, args, err := sq.
			Insert("Users").Columns("email", "password").
			Values(user.Email, string(hashed)).
			Suffix("RETURNING id, admin").
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			log.Debug(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var userID uint
		var admin bool

		err = dbClient.QueryRow(ctx, sql, args...).Scan(&userID, &admin)
		if err != nil {
			log.Debug(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		token, err := utils.GenerateToken(admin, userID)
		if err != nil {
			log.Debug(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			HTTPOnly: !c.IsFromLocal(),
			Secure:   !c.IsFromLocal(),
			MaxAge:   3600 * 24 * 7,
		})

		return c.JSON(fiber.Map{
			"token": token,
		})
	})

	// -------------------------------------------------------

	route.Post("/login", func(c *fiber.Ctx) error {
		authUser := &schemas.User{
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
		}

		if authUser.Email == "" || authUser.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email and password required",
			})
		}

		sql, args, err := sq.
			Select("id", "password", "admin").
			From("Users").
			Where(sq.Eq{"email": authUser.Email}).
			Limit(1).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			log.Debug(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var userID uint
		var hashedPassword string
		var admin bool

		err = dbClient.QueryRow(ctx, sql, args...).Scan(&userID, &hashedPassword, &admin)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid email or password",
				})
			}
			log.Fatal(err)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(authUser.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}

		token, err := utils.GenerateToken(admin, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			HTTPOnly: !c.IsFromLocal(),
			Secure:   !c.IsFromLocal(),
			MaxAge:   3600 * 24 * 7,
		})

		return c.JSON(fiber.Map{
			"token": token,
		})
	})
}
