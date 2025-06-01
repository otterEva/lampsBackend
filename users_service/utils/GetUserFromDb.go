package utils

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber"
	"github.com/jackc/pgx/v5"
	"github.com/otterEva/lamps/users_service/logs"
	"github.com/otterEva/lamps/users_service/settings"
)

func GetUserFromDb(c *fiber.Ctx, ctx context.Context, userId uint, admin bool) error {

	sql, args, err := sq.
		Select("1").
		From("Users").
		Where(sq.Eq{"id": userId, "admin": admin}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logs.Logger.Error("error while creating request", "sql", sql, "args", args)
		return c.SendStatus(StatusInternalServerError)
	}

	var exists int

	bdClient := settings.Clients.DbClient
	err = bdClient.QueryRow(ctx, sql, args...).Scan(&exists)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			logs.Logger.Debug("user not exists", "err", err.Error())

			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user doesn't exists",
		})
		}
		logs.Logger.Error("undexpected error")
		return c.SendStatus(StatusInternalServerError)
	}

	return nil
}
