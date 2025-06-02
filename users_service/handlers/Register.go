package handlers

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/otterEva/lamps/users_service/logs"
	"github.com/otterEva/lamps/users_service/schemas"
	"github.com/otterEva/lamps/users_service/settings"
	"github.com/otterEva/lamps/users_service/utils"
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

	logs.Logger.Debug("creds", "email", user.Email, "password", user.Password)

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logs.Logger.Debug("error while hashing password")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	sqlQuery, args, err := sq.
		Insert("Users").
		Columns("email", "password").
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

	err = dbClient.QueryRow(ctx, sqlQuery, args...).Scan(&userID, &admin)
	if err != nil {
		logs.Logger.Debug(
			"error while executing query",
			"sql", sqlQuery,
			"args", args,
			"err", err,
		)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	token, err := utils.GenerateToken(admin, userID)
	if err != nil {
		logs.Logger.Debug("creating token error", "err", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

c.Cookie(&fiber.Cookie{
    Name:     "jwt",
    Value:    token,
    HTTPOnly: false,         // JS должен иметь к ней доступ
    Secure:   true,          // в production — true, тут важно для SameSite=None
    SameSite: "None",        // без этого Lax по умолчанию и куку не пришлют по POST из другого origin
    Path:     "/",           // чтобы документ.cookie увидел куку на любом пути
    MaxAge:   3600 * 24 * 7, // 7 дней
})


	return nil
}
