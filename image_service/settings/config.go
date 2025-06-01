package settings

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type configStruct struct {
	APP_URL string `env:"APP_URL"`

	MinioRootUser     string `env:"MinioRootUser"`
	MinioRootPassword string `env:"MinioRootPassword"`
	MinioUseSSL       bool   `env:"MinioUseSSL"`
	MINIO_BUCKET      string `env:"MINIO_BUCKET"`

	MINIO_ACCESS_KEY string `env:"MINIO_ACCESS_KEY"`
	MINIO_SECRET_KEY string `env:"MINIO_SECRET_KEY"`
	MINIO_URL        string `env:"MINIO_URL"`
}

func getConfig() *configStruct {

	var cfg configStruct
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("failed to load config: %v", err.Error())
	}
	return &cfg
}

var Config *configStruct = getConfig()
