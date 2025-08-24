package postgre

import (
	"context"
	"gophkeeper/internal/storage/model"
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
