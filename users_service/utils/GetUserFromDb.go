package utils

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/otterEva/lamps/users_service/settings"
)

func GetUserFromDb(ctx context.Context, userId string, admin string) error {

	sql, args, err := sq.
		Select("1").
		From("Users").
		Where(sq.Eq{"id": userId, "admin": admin}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	var exists int

	bdClient := settings.Clients.DbClient
	err = bdClient.QueryRow(ctx, sql, args...).Scan(&exists)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user not found or access denied")
		}
		return fmt.Errorf("database error: %w", err)
	}

	return nil
}
