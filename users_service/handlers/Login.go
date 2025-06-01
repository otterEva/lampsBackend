package handlers

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/otterEva/lamps/users_service/schemas"
	"github.com/otterEva/lamps/users_service/settings"
	"github.com/otterEva/lamps/users_service/utils"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *fiber.Ctx, ctx context.Context) error {

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

	dbClient := settings.Clients.DbClient

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
	return nil
}
