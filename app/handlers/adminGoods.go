package handlers

import (
	"bytes"
	"context"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/otterEva/lamps/app/schemas"
	"github.com/otterEva/lamps/app/utils"
)

func AdminGoodsHandler(route fiber.Router, dbClient *pgxpool.Pool, ctx context.Context) {

	// CREATE -------------------------------------------------

	route.Post("/", func(c *fiber.Ctx) error {

		val := c.Locals("admin")
		isAdmin, ok := val.(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "No file provided")
		}

		file, err := fileHeader.Open()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Cannot open file")
		}
		defer file.Close()

		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(file); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to read file")
		}

		fileName, err := utils.AddFile(fileHeader.Filename, buf.Bytes(), fileHeader.Header.Get("Content-Type"))
		if err != nil {
			return err
		}

		// ---------------------------------------------------------------------------

		costStr := c.FormValue("cost")
		costUint64, err := strconv.ParseUint(costStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid cost",
			})
		}

		activeStr := c.FormValue("active")
		active, err := strconv.ParseBool(activeStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid active form",
			})
		}

		good := &schemas.Good{
			Description: c.FormValue("description"),
			Name:        c.FormValue("name"),
			ImageURL:    fileName,
			Cost:        uint(costUint64),
			Active:      active,
		}

		if good.Description == "" || good.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "description, name, and image_url are required",
			})
		}

		sql, args, err := sq.
			Insert("Goods").
			Columns("description", "name", "image_url", "cost", "active").
			Values(good.Description, good.Name, good.ImageURL, good.Cost, good.Active).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			log.Debug(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		_, err = dbClient.Exec(ctx, sql, args...)
		if err != nil {
			log.Debug(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.SendStatus(fiber.StatusCreated)
	})

	// UPDATE -------------------------------------------------

	route.Patch("/:id", func(c *fiber.Ctx) error {

		val := c.Locals("admin")
		isAdmin, ok := val.(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "No file provided")
		}

		file, err := fileHeader.Open()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Cannot open file")
		}
		defer file.Close()

		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(file); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to read file")
		}

		fileName, err := utils.AddFile(fileHeader.Filename, buf.Bytes(), fileHeader.Header.Get("Content-Type"))
		if err != nil {
			return err
		}

		idStr := c.Params("id")
		idUint64, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
		}
		id := uint(idUint64)

		update := sq.Update("Goods").Where(sq.Eq{"id": id})

		if desc := c.FormValue("description"); desc != "" {
			update = update.Set("description", desc)
		}
		if name := c.FormValue("name"); name != "" {
			update = update.Set("name", name)
		}
		if img := c.FormValue("image_url"); img != "" {
			update = update.Set(fileName, img)
		}
		if costStr := c.FormValue("cost"); costStr != "" {
			if costUint64, err := strconv.ParseUint(costStr, 10, 64); err == nil {
				update = update.Set("cost", uint(costUint64))
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cost"})
			}
		}
		if activeStr := c.FormValue("active"); activeStr != "" {
			if activeBool, err := strconv.ParseBool(activeStr); err == nil {
				update = update.Set("active", activeBool)
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid active"})
			}
		}

		sql, args, err := update.PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		_, err = dbClient.Exec(ctx, sql, args...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.SendStatus(fiber.StatusOK)
	})

	// DELETE -------------------------------------------------

	route.Delete("/:id", func(c *fiber.Ctx) error {

		val := c.Locals("admin")
		isAdmin, ok := val.(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}

		idStr := c.Params("id")
		idUint64, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
		}
		id := uint(idUint64)

		checkSql, checkArgs, err := sq.
			Select("1").
			From("Goods").
			Where(sq.Eq{"id": id}).
			Limit(1).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var exists int
		err = dbClient.QueryRow(ctx, checkSql, checkArgs...).Scan(&exists)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permission denied"})
		}

		sql, args, err := sq.
			Delete("Goods").
			Where(sq.Eq{"id": id}).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		cmdTag, err := dbClient.Exec(ctx, sql, args...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if cmdTag.RowsAffected() == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Good not found"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	// GET -------------------------------------------------

	route.Get("/", func(c *fiber.Ctx) error {

		val := c.Locals("admin")
		isAdmin, ok := val.(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}

		sql, args, err := sq.
			Select("id", "description", "name", "image_url", "cost", "active").
			From("Goods").
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		rows, err := dbClient.Query(ctx, sql, args...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		var goods []schemas.GoodDB

		for rows.Next() {
			var g schemas.GoodDB
			err := rows.Scan(&g.ID, &g.Description, &g.Name, &g.ImageURL, &g.Cost, &g.Active)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
			goods = append(goods, g)
		}

		return c.JSON(goods)
	})

}
