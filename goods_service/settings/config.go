package settings

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type configStruct struct {
	DATABASE_URL string `env:"DB_DSN"`

	SECRET  string `env:"SECRET"`
	APP_URL string `env:"APP_URL"`
	LOG_LEVEL string `env:"LOG_LEVEL"`
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
