package handlers

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/otterEva/lamps/users_service/logs"
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

	logs.Logger.Debug("user creds", "email", authUser.Email, "passwrod", authUser.Password)

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
		logs.Logger.Error("query error", "sql", sql, "args", args, "err", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var userID uint
	var hashedPassword string
	var admin bool

	dbClient := settings.Clients.DbClient

	err = dbClient.QueryRow(ctx, sql, args...).Scan(&userID, &hashedPassword, &admin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logs.Logger.Error("user doesnt exist")
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		logs.Logger.Error("unexpected db error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(authUser.Password)); err != nil {
		logs.Logger.Info("invalid credentials")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	logs.Logger.Debug("to token", "admin", admin, "userId", userID)

	token, err := utils.GenerateToken(admin, userID)
	if err != nil {
		logs.Logger.Error("error while creating a token")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

c.Cookie(&fiber.Cookie{
    Name:     "jwt",
    Value:    token,
    HTTPOnly: false,    
    Secure:   true,          
    SameSite: "None",        
    Path:     "/",           
    MaxAge:   3600 * 24 * 7, 
	Domain: "127.0.0.1",
})

	return nil
}
