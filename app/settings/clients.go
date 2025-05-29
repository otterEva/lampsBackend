package settings

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go"
)

type clientsStruct struct {
	DbClient    *pgxpool.Pool
	MinioClient *minio.Client
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

func getMinioClient() *minio.Client {

	minioClient, err := minio.New(Config.MINIO_URL, Config.MINIO_ACCESS_KEY, Config.MINIO_SECRET_KEY, false)
	if err != nil {
		log.Fatal(err)
	}

	err = minioClient.MakeBucket(Config.MINIO_BUCKET, "")
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(Config.MINIO_BUCKET)
		if errBucketExists == nil && exists {
			log.Debug("We already own %s\n", Config.MINIO_BUCKET)
		} else {
			log.Fatal(err)
		}
	} else {
		log.Debug("Successfully created %s\n", Config.MINIO_BUCKET)
	}

	return minioClient
}

func getClients() *clientsStruct {
	var clients clientsStruct
	clients.DbClient = getDbClient()
	clients.MinioClient = getMinioClient()
	return &clients
}

var Clients *clientsStruct = getClients()
