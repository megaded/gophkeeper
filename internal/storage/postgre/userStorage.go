// Содержит методы для работы с пользователем
package postgre

import (
	"context"
	"errors"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/storage/model"

	"gorm.io/gorm"
)

// Сохраняет нового пользователя в Бд
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

// Получение информации о пользователи по логину
func (s *PgStorage) GetUser(ctx context.Context, login string) (model.User, error) {
	var user model.User
	result := s.db.WithContext(ctx).Where("name = ?", login).First(&user)
	switch {
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return user, internal_error.ErrUserNotFound
	default:
		return user, result.Error
	}
}
