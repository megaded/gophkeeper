package minio

import (
	"context"
	"gophkeeper/internal/config"
	"gophkeeper/internal/logger"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

var (
	testBucker = "41622"
)

type MinioStorage struct {
	client *minio.Client
}

func NewStorage(cfg config.Config) (*MinioStorage, error) {
	client, err := minio.New(cfg.MinioAddress, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: false,
	})
	_, err = client.HealthCheck(time.Duration(time.Second * 1))
	if err != nil {
		return nil, err
	}

	ok, err := client.BucketExists(context.TODO(), testBucker)
	if !ok {
		err = client.MakeBucket(context.TODO(), testBucker, minio.MakeBucketOptions{})
		if err != nil {
			panic(err)
		}
	}
	logger.Log.Info("Minio UP")
	return &MinioStorage{client: client}, nil
}

func (m *MinioStorage) UploadFile(ctx context.Context, userName string, fileName string, reader io.Reader, size int64) error {
	_, err := m.client.PutObject(context.TODO(), testBucker, fileName, reader, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		logger.Log.Error("Minio. Загрузка файла", zap.Error(err))
	}
	return err
}
