package manager

import (
	"context"
	"gophkeeper/internal/config"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/server/dto"
	"gophkeeper/internal/storage/model"
)

type CreditCardManager struct {
	cryptoManager Crypter
	storager      creditCardStorager
}

func (c CreditCardManager) DeleteCreditCard(ctx context.Context, userId uint, id uint) error {
	card, err := c.storager.GetCreditCard(ctx, id)
	if err != nil {
		return err
	}
	if card.UserId != userId {
		return internal_error.ErrorAccessDenied
	}
	err = c.storager.DeleteCreditCard(ctx, id)
	return err
}

func (c CreditCardManager) AddCreditCard(ctx context.Context, userId uint, dto dto.Card) error {
	number, err := c.cryptoManager.Encrypt(dto.Number)
	if err != nil {
		return err
	}
	cvv, err := c.cryptoManager.Encrypt(dto.CVV)
	if err != nil {
		return err
	}
	exp, err := c.cryptoManager.Encrypt(dto.Exp)
	if err != nil {
		return err
	}
	return c.storager.AddCreditCard(ctx, userId, number, exp, cvv, dto.Description)
}

func (c CreditCardManager) GetCreditCards(ctx context.Context, userId uint) ([]dto.Card, error) {
	data, err := c.storager.GetCreditCards(ctx, userId)
	if err != nil {
		return nil, err
	}
	result := make([]dto.Card, 0, len(data))
	for _, card := range data {
		number, err := c.cryptoManager.Decrypt(card.Number)
		if err != nil {
			return nil, err
		}
		cvv, err := c.cryptoManager.Decrypt(card.CVE)
		if err != nil {
			return nil, err
		}
		exp, err := c.cryptoManager.Decrypt(card.Ext)
		if err != nil {
			return nil, err
		}
		result = append(result, dto.Card{Number: number, Exp: exp, CVV: cvv, Description: card.Description})
	}
	return result, nil
}

type creditCardStorager interface {
	AddCreditCard(ctx context.Context, userId uint, number []byte, ext []byte, cvv []byte, description string) error
	GetCreditCards(ctx context.Context, userId uint) ([]model.CreditCard, error)
	GetCreditCard(ctx context.Context, id uint) (model.CreditCard, error)
	DeleteCreditCard(ctx context.Context, id uint) error
}

func NewCreditCardManager(cfg config.Config, storager creditCardStorager) CreditCardManager {
	cryptoManager := NewCryptoManager(cfg)
	return CreditCardManager{cryptoManager: &cryptoManager, storager: storager}
}
