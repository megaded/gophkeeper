package storage

import (
	"context"
	"errors"
	"gophkeeper/internal/config"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/server/dto"
	"gophkeeper/internal/storage/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgStorage struct {
	db *gorm.DB
}

func (s PgStorage) AddUser(ctx context.Context, login string, password string) error {
	db := s.db.WithContext(ctx)
	db.Begin()
	defer db.Commit()
	var user model.User
	result := db.Where("name = ?", login).First(&user)
	switch {
	case result.Error == nil:
		return internal_error.ErrUserAlreadyExists
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		newUser := model.User{
			Name: login, Hash: password,
		}
		r := db.Create(&newUser)
		if r.Error != nil {
			db.Rollback()
			return r.Error
		}

		user = newUser
		return r.Error
	default:
		return result.Error
	}
}
func (s PgStorage) AddCredentials(ctx context.Context, dto dto.Credentials) error {
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
