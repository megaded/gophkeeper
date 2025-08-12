package filestorage

import (
	"gophkeeper/internal/config"

	"github.com/minio/minio-go"
)

type MinioStorage struct {
	client *minio.Client
}

func NewStorage(cfg config.Config) (*MinioStorage, error) {
	client, err := minio.New(cfg.MinioAddress, cfg.MinioAccessKey, cfg.MinioSecretKey,
		true)
	if err != nil {
		return nil, err
	}
	return &MinioStorage{client: client}, nil
}
