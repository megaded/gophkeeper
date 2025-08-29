package manager

import (
	"context"
	"gophkeeper/internal/dto"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/storage/model"
)

// Класс для работы с данными логин\пароль
type CredentialsManager struct {
	cryptoManager Crypter
	storage       credentialsStorager
}

func NewCredentialsManager(cryptoManager Crypter, storage credentialsStorager) CredentialsManager {
	return CredentialsManager{cryptoManager: cryptoManager, storage: storage}
}

// Добавляет тип данных логин\пароль для пользователя
func (c CredentialsManager) AddCredentials(ctx context.Context, userId uint, cred dto.Credentials) error {
	login, err := c.cryptoManager.Encrypt(cred.Login)
	if err != nil {
		return err
	}
	password, err := c.cryptoManager.Encrypt(cred.Password)
	if err != nil {
		return err
	}
	return c.storage.AddCredentials(ctx, userId, login, password, cred.Description)
}

// Получение списка данных логин\пароль по Id пользователя
func (c CredentialsManager) GetCredentials(ctx context.Context, userId uint) ([]dto.Credentials, error) {
	creds, err := c.storage.GetCredentials(ctx, userId)
	if err != nil {
		return nil, err
	}
	result := make([]dto.Credentials, 0, len(creds))
	for _, cred := range creds {
		login, err := c.cryptoManager.Decrypt(cred.Login)
		if err != nil {
			return nil, err
		}
		password, err := c.cryptoManager.Decrypt(cred.Password)
		if err != nil {
			return nil, err
		}
		result = append(result, dto.Credentials{Login: login, Password: password, Description: cred.Description, Id: cred.ID})
	}
	return result, err
}

// Удаляет тип данных логин\пароль по ID
func (c CredentialsManager) DeleteCredential(ctx context.Context, userId uint, id uint) error {
	cred, err := c.storage.GetCredential(ctx, id)
	if err != nil {
		return err
	}
	if cred.UserId != userId {
		return internal_error.ErrorAccessDenied
	}
	err = c.storage.DeleteCredentials(ctx, id)
	return err
}

// Обновление данных логин\пароль
func (c CredentialsManager) UpdateCredentials(ctx context.Context, userId uint, dto dto.Credentials) error {
	cred, err := c.storage.GetCredential(ctx, dto.Id)
	if err != nil {
		return err
	}
	if cred.UserId != userId {
		return internal_error.ErrorAccessDenied
	}
	login, err := c.cryptoManager.Encrypt(dto.Login)
	if err != nil {
		return err
	}
	password, err := c.cryptoManager.Encrypt(dto.Password)
	if err != nil {
		return err
	}
	return c.storage.UpdateCredentials(ctx, userId, login, password, cred.Description)
}

type credentialsStorager interface {
	AddCredentials(ctx context.Context, userId uint, login []byte, password []byte, description string) error
	GetCredentials(ctx context.Context, userId uint) ([]model.Credentials, error)
	GetCredential(ctx context.Context, id uint) (model.Credentials, error)
	DeleteCredentials(ctx context.Context, id uint) error
	UpdateCredentials(ctx context.Context, id uint, login []byte, password []byte, description string) error
}
