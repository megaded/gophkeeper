package manager

import (
	"context"
	"gophkeeper/internal/server/dto"
	"io"
	"strconv"
)

type BinaryManager struct {
	fileStorager fileStorager
	storager     fileMetaStorager
}

func (b BinaryManager) UploadFile(ctx context.Context, userId uint, dto dto.BinaryFile, reader io.Reader) error {
	name, err := b.fileStorager.UploadFile(ctx, strconv.Itoa(int(userId)), dto.FileName, reader)
	if err != nil {
		return err
	}
	err = b.storager.AddBinary(ctx, userId, dto.Description, dto.FileName, name)
	if err != nil {
		b.fileStorager.DeleteFile(ctx, userId, name)
		return err
	}
	return nil
}

type fileStorager interface {
	UploadFile(ctx context.Context, userId string, fileName string, reader io.Reader) (string, error)
	DeleteFile(ctx context.Context, userId uint, name string)
}

type fileMetaStorager interface {
	AddBinary(ctx context.Context, userId uint, description string, orininalFileName string, externalFileName string) error
}
