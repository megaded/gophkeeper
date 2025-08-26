package manager

import (
	"context"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/server/dto"
	"gophkeeper/internal/storage/model"
)

type TextManager struct {
	storager textStorager
}

func NewTextManager(storager textStorager) TextManager {
	return TextManager{storager: storager}
}

type textStorager interface {
	AddText(ctx context.Context, userId uint, content string, description string) error
	GetTextList(ctx context.Context, userId uint) ([]model.Text, error)
	GetText(ctx context.Context, id uint) (model.Text, error)
	UpdateText(ctx context.Context, id uint, content string, description string) error
}

func (f TextManager) UploadText(ctx context.Context, dto dto.Text) error {
	return f.storager.AddText(ctx, dto.UserId, dto.Content, dto.Description)
}

func (f TextManager) GetTextList(ctx context.Context, userId uint) ([]dto.Text, error) {
	data, err := f.storager.GetTextList(ctx, userId)
	if err != nil {
		return nil, err
	}
	r := make([]dto.Text, 0, len(data))
	for _, i := range data {
		r = append(r, dto.Text{Id: i.ID, IsFile: i.IsFile, Description: i.Description})
	}
	return r, nil
}

// Обновление текстовые данные
func (c TextManager) UpdateText(ctx context.Context, userId uint, dto dto.Text) error {
	text, err := c.storager.GetText(ctx, dto.Id)
	if err != nil {
		return err
	}
	if text.UserId != userId {
		return internal_error.ErrorAccessDenied
	}
	return c.storager.UpdateText(ctx, userId, dto.Content, dto.Description)
}
