package postgre

import (
	"gophkeeper/internal/config"
	"gophkeeper/internal/storage/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgStorage struct {
	db *gorm.DB
}

// Создание новое хранилище PosgreSql
func NewStorage(c *config.Config) (PgStorage, error) {
	db, err := gorm.Open(postgres.Open(c.DBConnString), &gorm.Config{})
	if err != nil {
		return PgStorage{}, err
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Credentials{})
	db.AutoMigrate(&model.CreditCard{})
	db.AutoMigrate(&model.Binary{})
	db.AutoMigrate(&model.Text{})
	return PgStorage{db: db}, nil
}
