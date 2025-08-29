// Содержит методы для работы с типом текст
package postgre

import (
	"context"
	"errors"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/storage/model"

	"github.com/guregu/null/v6"
	"gorm.io/gorm"
)

func (s PgStorage) AddBinary(ctx context.Context, userId uint, description string, originalFileName string, externalFileName string) (model.Binary, error) {
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
		return model, r.Error
	}
	return model, nil
}

func (s PgStorage) AddTextFile(ctx context.Context, userId uint, description string, binaryId uint) error {
	db := s.db.WithContext(ctx)
	db.Begin()
	defer db.Commit()
	model := model.Text{}
	model.UserId = userId
	model.Description = description
	model.BinaryId = null.Int32From(int32(binaryId))
	model.IsFile = true
	r := db.Create(&model)
	if r.Error != nil {
		db.Rollback()
		return r.Error
	}
	return nil
}
func (s PgStorage) GetFileInfo(ctx context.Context, fileId uint) (model.Binary, error) {
	var model model.Binary
	result := s.db.WithContext(ctx).Where("id = ?", fileId).First(&model)
	switch {
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return model, internal_error.ErrRecordNotFound
	default:
		return model, result.Error
	}
}

// Возвращает список данных логин\пароль по Id юзера
func (s PgStorage) GetBinaryFiles(ctx context.Context, userId uint) ([]model.Binary, error) {
	var model []model.Binary
	result := s.db.WithContext(ctx).Where("user_id = ?", userId).Find(&model)
	switch {
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return nil, internal_error.ErrRecordNotFound
	default:
		return model, result.Error
	}
}

// Обновляет бинарный файл
func (s PgStorage) UpdateBinary(ctx context.Context, id uint, externalName string, originalName string, description string) error {
	r := s.db.WithContext(ctx).Model(&model.Binary{}).Where("id = ?", id).Updates(model.Binary{ExternalFileName: externalName, OriginalFileName: originalName, Description: description})
	if r.Error != nil {
		return r.Error
	}
	return nil
}

// Удаление бинарных данных по Id
func (s PgStorage) DeleteBinary(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&model.Binary{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
