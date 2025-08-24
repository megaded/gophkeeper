// Содержит методы для работы с типом данных логин\пароль
package postgre

import (
	"context"
	"errors"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/server/dto"
	"gophkeeper/internal/storage/model"

	"gorm.io/gorm"
)

// Создает новую запись в БД типа логин\пароль
func (s PgStorage) AddCredentials(ctx context.Context, userId uint, login []byte, password []byte, description string) error {
	db := s.db.WithContext(ctx)
	db.Begin()
	defer db.Commit()
	newCred := model.Credentials{
		Login: login, Password: password,
	}
	newCred.UserId = userId
	newCred.Description = description
	r := db.Create(&newCred)
	if r.Error != nil {
		db.Rollback()
		return r.Error
	}
	return nil
}

// Возвращает  данных логин\пароль по Id
func (s PgStorage) GetCredential(ctx context.Context, id uint) (model.Credentials, error) {
	var model model.Credentials
	result := s.db.WithContext(ctx).Where("id = ?", id).First(&model)
	switch {
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return model, internal_error.ErrUserNotFound
	default:
		return model, result.Error
	}
}

// Возвращает список данных логин\пароль по Id юзера
func (s PgStorage) GetCredentials(ctx context.Context, userId uint) ([]model.Credentials, error) {
	var model []model.Credentials
	result := s.db.WithContext(ctx).Where("user_id = ?", userId).Find(&model)
	switch {
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return nil, internal_error.ErrRecordNotFound
	default:
		return model, result.Error
	}
}

// Удаление данных логин\пароль по Id
func (s PgStorage) DeleteCredentials(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(model.Credentials{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Обновляет данные типа логин\пароль
func (s PgStorage) UpdateCredentials(ctx context.Context, cred dto.Credentials) error {
	return nil
}
