package settings

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type clientsStruct struct {
	DbClient *pgxpool.Pool
}

func getDbClient() *pgxpool.Pool {

	ctx := context.Background()
	pool, err := pgxpool.New(context.Background(), Config.DATABASE_URL)

	if err != nil {
		log.Debug(err.Error())
		log.Fatal("Unable to connect to database:", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Debug(err.Error())
		log.Fatal("Unable to ping database:", err)
	}

	return pool
}

func getClients() *clientsStruct {
	var clients clientsStruct
	clients.DbClient = getDbClient()
	return &clients
}

var Clients *clientsStruct = getClients()
