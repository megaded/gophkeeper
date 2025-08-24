// Содержит методы для работы с типом текст
package postgre

import (
	"context"
	"gophkeeper/internal/storage/model"
)

func (s PgStorage) AddBinary(ctx context.Context, userId uint, description string, originalFileName string, externalFileName string) (uint, error) {
	db := s.db.WithContext(ctx)
	db.Begin()
	defer db.Commit()
	model := model.Binary{}
	model.UserId = userId
	model.Description = description
	model.OriginalFileName = originalFileName
	model.ExternalFileName = externalFileName
	r := db.Create(&model)
	if r.Error != nil {
		db.Rollback()
		return 0, r.Error
	}
	return model.ID, nil
}

func (s PgStorage) AddTextFile(ctx context.Context, userId uint, description string, binaryId uint) error {
	db := s.db.WithContext(ctx)
	db.Begin()
	defer db.Commit()
	model := model.Text{}
	model.UserId = userId
	model.Description = description
	model.BinaryId = binaryId
	model.IsFile = true
	r := db.Create(&model)
	if r.Error != nil {
		db.Rollback()
		return r.Error
	}
	return nil
}
