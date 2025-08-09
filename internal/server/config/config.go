package config

import "github.com/caarlos0/env"

type Config struct {
	DBConnString   string `env:"DATABASE_URI"`
	MinioAddress   string `env:"MINIO_ADDRESS"`
	MinioAccessKey string `env:"MINIO_ACCESS_KEY"`
}

func GetConfig() Config {
	cfg := Config{}
	env.Parse(cfg)
	return cfg
}
