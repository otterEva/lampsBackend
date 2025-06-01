package handlers

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/otterEva/lamps/users_service/logs"
	"github.com/otterEva/lamps/users_service/schemas"
	"github.com/otterEva/lamps/users_service/settings"
	"github.com/otterEva/lamps/users_service/utils"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *fiber.Ctx, ctx context.Context) error {

	user := &schemas.User{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	if user.Email == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password required",
		})
	}

	logs.Logger.Debug("creds", "email", email, "password", password)
			
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logs.Logger.Debug("error while hashing password")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	sql, args, err := sq.
		Insert("Users").Columns("email", "password").
		Values(user.Email, string(hashed)).
		Suffix("RETURNING id, admin").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logs.Logger.Error("error while creating query")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var userID uint
	var admin bool

	dbClient := settings.Clients.DbClient

	err = dbClient.QueryRow(ctx, sql, args...).Scan(&userID, &admin)
	if err != nil {
		logs.Logger.Debug("error while executing query", "sql", sql, "args", args, "err": err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	token, err := utils.GenerateToken(admin, userID)

	if err != nil {
		log.Debug(err.Error())
		logs.Logger.Debug("creating token error")
		return c.SendStatus(fiber.StatusInternalServerError)
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
