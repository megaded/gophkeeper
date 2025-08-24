package minio

import (
	"context"
	"gophkeeper/internal/config"
	"gophkeeper/internal/logger"
	"io"
	"path/filepath"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

var (
	bucketPrefix = "keeper"
)

type MinioStorage struct {
	client *minio.Client
}

// Создание новое файловое хранилище на основе клиента Minio
func NewStorage(cfg config.Config) (*MinioStorage, error) {
	client, err := minio.New(cfg.MinioAddress, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: false,
	})
	_, err = client.HealthCheck(time.Duration(time.Second * 1))
	if err != nil {
		return nil, err
	}

	ok, err := client.BucketExists(context.TODO(), bucketPrefix)
	if !ok {
		err = client.MakeBucket(context.TODO(), bucketPrefix, minio.MakeBucketOptions{})
		if err != nil {
			panic(err)
		}
	}
	logger.Log.Info("Minio UP")
	return &MinioStorage{client: client}, nil
}

// Загружает файл в файловое хранилище.
// Возвращает внешнее имя файла и ошибку
func (m *MinioStorage) UploadFile(ctx context.Context, userId string, fileName string, reader io.Reader) (string, error) {
	externalName := generateFileName(fileName)
	bucketName := getBucketName(userId)
	ok, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return "", err
	}
	if !ok {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", err
		}
	}
	_, err = m.client.PutObject(context.TODO(), bucketPrefix, externalName, reader, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		logger.Log.Error("Minio. Загрузка файла", zap.Error(err))
		return "", err
	}
	return externalName, err
}

func (m *MinioStorage) DownloadFile(ctx context.Context, userId uint, fileName string) (io.Reader, error) {
	bucketName := getBucketName(strconv.Itoa(int(userId)))
	r, err := m.client.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (m MinioStorage) DeleteFile(ctx context.Context, userId uint, name string) {
	panic(" MinioStorage DeleteFile")
}

func getBucketName(userId string) string {
	return bucketPrefix + string(userId)
}

func generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	return ext
}
