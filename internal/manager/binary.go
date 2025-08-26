// Логикак по работе с бинарными файлами
package manager

import (
	"context"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/server/dto"
	"gophkeeper/internal/storage/model"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type BinaryManager struct {
	fileStorager fileStorager
	storager     fileMetaStorager
}

// Создание менеджера по работе с файлами
func NewBinaryManager(fileStorager fileStorager, storager fileMetaStorager) BinaryManager {
	return BinaryManager{fileStorager: fileStorager, storager: storager}
}

func (b BinaryManager) uploadInternal(ctx context.Context, userId uint, dtos dto.BinaryFile, reader io.Reader) (dto.BinaryFile, error) {
	name, err := b.fileStorager.UploadFile(ctx, strconv.Itoa(int(userId)), dtos.FileName, reader)
	if err != nil {
		return dto.BinaryFile{}, err
	}
	new, err := b.storager.AddBinary(ctx, userId, dtos.Description, dtos.FileName, name)
	if err != nil {
		b.fileStorager.DeleteFile(ctx, userId, name)
		return dto.BinaryFile{}, err
	}
	return dto.BinaryFile{Id: new.ID, UserId: new.ID, Description: new.Description, FileName: new.OriginalFileName, ExternalFileName: new.ExternalFileName}, nil
}

// Загрузка файла
func (b BinaryManager) UploadFile(ctx context.Context, userId uint, dto dto.BinaryFile, reader io.Reader) (dto.BinaryFile, error) {
	return b.uploadInternal(ctx, userId, dto, reader)
}

// Загрузка текстового файла
func (b BinaryManager) UploadTextFile(ctx context.Context, userId uint, dto dto.BinaryFile, reader io.Reader) error {
	name, err := b.fileStorager.UploadFile(ctx, strconv.Itoa(int(userId)), dto.FileName, reader)
	if err != nil {
		return err
	}
	newFile, err := b.storager.AddBinary(ctx, userId, dto.Description, dto.FileName, name)
	if err != nil {
		b.fileStorager.DeleteFile(ctx, userId, name)
		return err
	}
	err = b.storager.AddTextFile(ctx, userId, dto.Description, newFile.ID)
	if err != nil {
		b.fileStorager.DeleteFile(ctx, userId, name)
		return err
	}
	return nil
}

// Скачивание файла
func (b BinaryManager) DownloadFile(ctx context.Context, userId uint, id uint) (reader io.Reader, info dto.BinaryFile, err error) {
	fileInfo, err := b.storager.GetFileInfo(ctx, id)
	if err != nil {
		return nil, dto.BinaryFile{}, nil
	}
	if fileInfo.UserId != userId {
		return nil, dto.BinaryFile{}, internal_error.ErrorAccessDenied
	}
	r, err := b.fileStorager.DownloadFile(ctx, userId, fileInfo.ExternalFileName)
	if err != nil {
		return nil, dto.BinaryFile{}, err
	}
	//SaveLocalFile(r, fileInfo.OriginalFileName)
	return r, dto.BinaryFile{Id: fileInfo.ID, FileName: fileInfo.OriginalFileName}, nil

}

// Получение списка загруженных файлов пользователя
func (b BinaryManager) GetBinaryFiles(ctx context.Context, userId uint) ([]dto.BinaryFile, error) {
	files, err := b.storager.GetBinaryFiles(ctx, userId)
	if err != nil {
		return nil, nil
	}
	result := make([]dto.BinaryFile, 0, len(files))
	for _, f := range files {
		result = append(result, dto.BinaryFile{Id: f.ID, FileName: f.OriginalFileName, Description: f.Description})
	}
	return result, err
}

// Получение списка загруженных файлов пользователя
func (b BinaryManager) UpdateBinaryFile(ctx context.Context, userId uint, dto dto.BinaryFile, reader io.Reader) error {
	fileInfo, err := b.storager.GetFileInfo(ctx, dto.Id)
	if err != nil {
		return err
	}
	if fileInfo.UserId != userId {
		return internal_error.ErrorAccessDenied
	}
	newFile, err := b.uploadInternal(ctx, userId, dto, reader)
	if err != nil {
		return err
	}
	err = b.storager.UpdateBinary(ctx, dto.Id, newFile.ExternalFileName, dto.FileName, dto.Description)
	if err != nil {
		return err
	}
	err = b.fileStorager.DeleteFile(ctx, userId, fileInfo.ExternalFileName)
	if err != nil {
		return err
	}
	return nil
}

// Удаляет бинарный файл по Id
func (b BinaryManager) DeleteBinaryFile(ctx context.Context, userId uint, id uint) error {
	fileInfo, err := b.storager.GetFileInfo(ctx, id)
	if err != nil {
		return err
	}
	if fileInfo.UserId != userId {
		return internal_error.ErrorAccessDenied
	}
	return nil
}

func SaveLocalFile(reader io.Reader, fileName string) error {
	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}
	downloadDir := filepath.Join(rootDir, "Download")
	_, err = os.Stat(downloadDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(downloadDir, os.ModeAppend)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(filepath.Join(downloadDir, fileName))
	if err != nil {
		return nil
	}
	defer file.Close()
	_, err = io.Copy(file, reader)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

type fileStorager interface {
	UploadFile(ctx context.Context, userId string, fileName string, reader io.Reader) (string, error)
	DeleteFile(ctx context.Context, userId uint, name string) error
	DownloadFile(ctx context.Context, userId uint, fileName string) (io.Reader, error)
}

type fileMetaStorager interface {
	AddBinary(ctx context.Context, userId uint, description string, originalFileName string, externalFileName string) (model.Binary, error)
	AddTextFile(ctx context.Context, userId uint, description string, binaryId uint) error
	GetFileInfo(ctx context.Context, fileId uint) (model.Binary, error)
	GetBinaryFiles(ctx context.Context, userId uint) ([]model.Binary, error)
	UpdateBinary(ctx context.Context, id uint, externalName string, originalName string, description string) error
}
