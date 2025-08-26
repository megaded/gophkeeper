package postgre

import (
	"context"
	"errors"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/storage/model"

	"gorm.io/gorm"
)

func (s PgStorage) AddText(ctx context.Context, userId uint, content string, description string) error {
	db := s.db.WithContext(ctx)
	db.Begin()
	defer db.Commit()
	model := model.Text{}
	model.UserId = userId
	model.Description = description
	model.IsFile = false
	model.Content = content
	r := db.Create(&model)
	if r.Error != nil {
		db.Rollback()
		return r.Error
	}
	return nil
}

// Возвращает текстовые данные по Id
func (s PgStorage) GetText(ctx context.Context, id uint) (model.Text, error) {
	var model model.Text
	result := s.db.WithContext(ctx).Where("id = ?", id).First(&model)
	switch {
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return model, internal_error.ErrUserNotFound
	default:
		return model, result.Error
	}
}

// Возвращает список текстовых данных пользователя по Id
func (s PgStorage) GetTextList(ctx context.Context, userId uint) ([]model.Text, error) {
	var model []model.Text
	result := s.db.WithContext(ctx).Where("user_id = ?", userId).Find(&model)
	switch {
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return nil, internal_error.ErrRecordNotFound
	default:
		return model, result.Error
	}
}

// Обновление текстовых данных
func (s PgStorage) UpdateText(ctx context.Context, id uint, content string, description string) error {
	r := s.db.Model(&model.Text{}).Where("id = ?", id).Updates(model.Text{Content: content, Description: description})
	if r.Error != nil {
		return r.Error
	}
	return nil
}

// Удаление текстовый данных
func (s PgStorage) DeleteText(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(model.Text{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
