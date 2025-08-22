package postgre

import (
	"context"
	"gophkeeper/internal/server/dto"
	"gophkeeper/internal/storage/model"
)

func (s PgStorage) AddCredentials(ctx context.Context, userId uint, login []byte, password []byte) error {
	db := s.db.WithContext(ctx)
	db.Begin()
	defer db.Commit()
	newCred := model.Credentials{
		Login: login, Password: password,
	}
	newCred.UserId = userId
	r := db.Create(&newCred)
	if r.Error != nil {
		db.Rollback()
		return r.Error
	}
	return nil
}
func (s PgStorage) GetCredentials(ctx context.Context, userId uint) error {
	return nil
}
func (s PgStorage) DeleteCredentials(ctx context.Context, id uint) error {
	return nil
}
func (s PgStorage) UpdateCredentials(ctx context.Context, cred dto.Credentials) error {
	return nil
}
