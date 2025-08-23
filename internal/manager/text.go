package manager

import (
	"context"
	"gophkeeper/internal/server/dto"
)

type TextManager struct {
	storager textStorager
}

type textStorager interface {
	AddText(ctx context.Context, userId uint, content string, description string) error
}

func (f TextManager) UploadText(ctx context.Context, dto dto.Text) error {
	return f.storager.AddText(ctx, dto.UserId, dto.Content, dto.Description)
}
