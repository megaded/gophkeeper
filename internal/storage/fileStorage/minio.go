package filestorage

import (
	"gophkeeper/internal/server/config"

	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
)

type MinioStorage struct {
	client *minio.Client
}

func NewStorage(cfg config.Config) (MinioStorage, error) {
	client, err := minio.New(cf, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
