package config

import "github.com/caarlos0/env"

type Config struct {
	DBConnString   string `env:"DATABASE_URI,required"`
	MinioAddress   string `env:"MINIO_ADDRESS,required"`
	MinioAccessKey string `env:"MINIO_ACCESS_KEY,required"`
	MinioSecretKey string `env:"MINIO_SECRET_KEY,required"`
	SecretKey      string `env:"SECRET_KEY,required"`
}

func GetConfig() Config {
	cfg := Config{}
	env.Parse(cfg)
	return cfg
}
