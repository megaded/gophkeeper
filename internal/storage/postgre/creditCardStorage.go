// Содержит методы для работы с типом данных банковская карта
package postgre

import (
	"context"
	"errors"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/storage/model"

	"gorm.io/gorm"
)

// Удаляет тип бансковская карта по Id
func (s PgStorage) DeleteCreditCard(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(model.CreditCard{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Получает список банковских карт по Id пользователя
func (s PgStorage) GetCreditCards(ctx context.Context, userId uint) ([]model.CreditCard, error) {
	var cards []model.CreditCard
	result := s.db.WithContext(ctx).Where("user_id = ?", userId).Find(&cards)
	switch {
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return nil, internal_error.ErrRecordNotFound
	default:
		return cards, result.Error
	}
}

// Возвращает банковскую карту по ID
func (s PgStorage) GetCreditCard(ctx context.Context, id uint) (model.CreditCard, error) {
	var model model.CreditCard
	result := s.db.WithContext(ctx).Where("id = ?", id).First(&model)
	switch {
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return model, internal_error.ErrUserNotFound
	default:
		return model, result.Error
	}
}

// Добавляет в Бд тип данных банковская карта
func (s PgStorage) AddCreditCard(ctx context.Context, userId uint, number []byte, ext []byte, cvv []byte, description string) error {
	db := s.db.WithContext(ctx)
	db.Begin()
	defer db.Commit()
	card := model.CreditCard{
		Number: number,
		CVE:    cvv,
		Ext:    ext,
	}
	card.UserId = userId
	card.Description = description
	r := db.Create(&card)
	if r.Error != nil {
		db.Rollback()
		return r.Error
	}
	return nil
}

// Обновляет данные типа логин\пароль
func (s PgStorage) UpdateCreditCard(ctx context.Context, id uint, cvv []byte, exp []byte, cve []byte, description string) error {
	r := s.db.Model(&model.CreditCard{}).Where("id = ?", id).Updates(model.CreditCard{Number: cvv, Ext: exp, CVE: cve, Description: description})
	if r.Error != nil {
		return r.Error
	}
	return nil
}
