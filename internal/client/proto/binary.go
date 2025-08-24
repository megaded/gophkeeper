package proto

import (
	"bufio"
	"context"
	"gophkeeper/internal/logger"
	pb "gophkeeper/proto"
	"io"

	"go.uber.org/zap"
)

func (c *keeperClient) UploadBinaryFile(reader io.Reader, fileName string, description string, size int64) error {
	ctx, err := getCtx(c.token)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	k, err := c.client.UploadBinaryFile(ctx)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	buf := bufio.NewReader(reader)
	data := make([]byte, buf.Size())
	var totalSize int64 = 0
	for {
		b, err := buf.Read(data)
		totalSize = totalSize + int64(b)
		if err == io.EOF {
			logger.Log.Info("байтов", zap.Int64("всего", totalSize))
			logger.Log.Info("Конец файла")
			break
		}
		err = k.Send(&pb.UploadBinaryFileRequest{Content: data, Filename: fileName, Description: description, Size: size})
		if err != nil && err != io.EOF {
			_, err = k.CloseAndRecv()
			logger.Log.Error(err.Error())
			return err
		}
	}
	_, err = k.CloseAndRecv()
	logger.Log.Info("Закончили отправку")
	if err != nil && err != io.EOF {
		logger.Log.Error(err.Error())
	}
	return err
}

func (c *keeperClient) DownloadBinaryFile(ctx context.Context, id uint) (io.Reader, error) {
	r, err := c.client.DownloadBinaryFile(ctx, &pb.DownloadBinaryFileRequest{Id: int64(id)})
	if err != nil {
		return nil, nil
	}
	return r.Header(), nil
}
