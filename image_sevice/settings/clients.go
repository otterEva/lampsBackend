package settings

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/minio/minio-go"
)

type clientsStruct struct {
	MinioClient *minio.Client
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
	clients.MinioClient = getMinioClient()
	return &clients
}

var Clients *clientsStruct = getClients()
