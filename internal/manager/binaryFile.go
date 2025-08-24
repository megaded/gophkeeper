package manager

import (
	"context"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/server/dto"
	"gophkeeper/internal/storage/model"
	"io"
	"strconv"
)

type BinaryManager struct {
	fileStorager fileStorager
	storager     fileMetaStorager
}

func NewBinaryManager(fileStorager fileStorager, storager fileMetaStorager) BinaryManager {
	return BinaryManager{fileStorager: fileStorager, storager: storager}
}

func (b BinaryManager) UploadFile(ctx context.Context, userId uint, dto dto.BinaryFile, reader io.Reader) error {
	name, err := b.fileStorager.UploadFile(ctx, strconv.Itoa(int(userId)), dto.FileName, reader)
	if err != nil {
		return err
	}
	_, err = b.storager.AddBinary(ctx, userId, dto.Description, dto.FileName, name)
	if err != nil {
		b.fileStorager.DeleteFile(ctx, userId, name)
		return err
	}
	return nil
}

func (b BinaryManager) UploadTextFile(ctx context.Context, userId uint, dto dto.BinaryFile, reader io.Reader) error {
	name, err := b.fileStorager.UploadFile(ctx, strconv.Itoa(int(userId)), dto.FileName, reader)
	if err != nil {
		return err
	}
	id, err := b.storager.AddBinary(ctx, userId, dto.Description, dto.FileName, name)
	if err != nil {
		b.fileStorager.DeleteFile(ctx, userId, name)
		return err
	}
	err = b.storager.AddTextFile(ctx, userId, dto.Description, id)
	if err != nil {
		b.fileStorager.DeleteFile(ctx, userId, name)
		return err
	}
	return nil
}

func (b BinaryManager) DownloadFile(ctx context.Context, userId uint, id uint) (reader io.Reader, info dto.BinaryFile, err error) {
	fileInfo, err := b.storager.GetFileInfo(ctx, id)
	if err != nil {
		return nil, dto.BinaryFile{}, nil
	}
	if fileInfo.UserId != userId {
		return nil, dto.BinaryFile{}, internal_error.ErrorAccessDenied
	}
	r, err := b.fileStorager.DownloadFile(ctx, userId, fileInfo.ExternalFileName)
	return r, dto.BinaryFile{Id: fileInfo.ID, FileName: fileInfo.OriginalFileName}, nil

}

type fileStorager interface {
	UploadFile(ctx context.Context, userId string, fileName string, reader io.Reader) (string, error)
	DeleteFile(ctx context.Context, userId uint, name string)
	DownloadFile(ctx context.Context, userId uint, fileName string) (io.Reader, error)
}

type fileMetaStorager interface {
	AddBinary(ctx context.Context, userId uint, description string, originalFileName string, externalFileName string) (uint, error)
	AddTextFile(ctx context.Context, userId uint, description string, binaryId uint) error
	GetFileInfo(ctx context.Context, fileId uint) (model.Binary, error)
}
