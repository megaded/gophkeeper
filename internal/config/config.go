package config

import (
	"gophkeeper/internal/logger"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DBConnString   string `env:"DATABASE_URI,required"`
	MinioAddress   string `env:"MINIO_ADDRESS,required"`
	MinioAccessKey string `env:"MINIO_ACCESS_KEY,required"`
	MinioSecretKey string `env:"MINIO_SECRET_KEY,required"`
	SecretKey      string `env:"SECRET_KEY,required"`
	Address        string `env:"ADDRESS,required"`
}

func GetConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	cfg := Config{}
	env.Parse(&cfg)
	logger.Log.Info(cfg.Address)
	return cfg
}
