package postgre

import (
	"context"
	"errors"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/storage/model"

	"gorm.io/gorm"
)

func (s PgStorage) DeleteCreditCard(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(model.CreditCard{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
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
