package postgre

import (
	"gophkeeper/internal/config"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/storage/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgStorage struct {
	db *gorm.DB
}

// Создание новое хранилище PosgreSql
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
