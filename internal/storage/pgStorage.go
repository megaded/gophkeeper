package storage

import (
	"context"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/server/config"
	"gophkeeper/internal/server/dto"
	"gophkeeper/internal/storage/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgStorage struct {
	db *gorm.DB
}

func (s PgStorage) AddUser(ctx context.Context, login string, password string) error {
	return nil
}
func (s PgStorage) AddCredentials(ctx context.Context, login string, password string) error {
	return nil
}
func (s PgStorage) GetCredentials(ctx context.Context, login string) error {
	return nil
}
func (s PgStorage) DeleteCredentials(ctx context.Context, login string) error {
	return nil
}
func (s PgStorage) UpdateCredentials(ctx context.Context, cred dto.Credentials) error {
	return nil
}

func NewStorage(c *config.Config) PgStorage {
	db, err := gorm.Open(postgres.Open(c.DBConnString), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Credentials{})
	db.AutoMigrate(&model.CreditCard{})
	db.AutoMigrate(&model.Binary{})
	db.AutoMigrate(&model.Text{})
	return PgStorage{db: db}
}
