package configs

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBHost     string `envconfig:"DB_HOST"`
	DBPort     int    `envconfig:"DB_PORT"`
	DBUser     string `envconfig:"DB_USER"`
	DBPassword string `envconfig:"DB_PASSWORD"`
	DBName     string `envconfig:"DB_NAME"`
}

var Env Config

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	if err := godotenv.Load(); err != nil {
		return err
	}

	if err := envconfig.Process("", &Env); err != nil {
		return err
	}

	return nil
}
