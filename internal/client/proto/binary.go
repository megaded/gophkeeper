package proto

import (
	"bufio"
	"context"
	fileutil "gophkeeper/internal/client/proto/fileUtil"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/server/dto"
	pb "gophkeeper/proto"
	"io"

	"go.uber.org/zap"
)

func (c *keeperClient) UploadBinaryFile(reader io.Reader, fileName string, description string) error {
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
	for err != io.EOF {
		b, err := buf.Read(data)
		totalSize = totalSize + int64(b)
		if err == io.EOF {
			logger.Log.Info("байтов", zap.Int64("всего", totalSize))
			logger.Log.Info("Конец файла")
			break
		}
		err = k.Send(&pb.UploadBinaryFileRequest{Content: data, Filename: fileName, Description: description})
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

func (c *keeperClient) DownloadBinaryFile(ctx context.Context, id uint) error {
	ctx, err := getCtx(c.token)
	if err != nil {
		return err
	}
	r, err := c.client.DownloadBinaryFile(ctx, &pb.DownloadBinaryFileRequest{Id: int64(id)})
	if err != nil {
		return nil
	}
	rd, wr := io.Pipe()
	resp, err := r.Recv()
	if err != nil {
		return err
	}

	fineName := resp.Filename
	go func() {
		defer wr.Close()

		wr.Write(resp.Content)
		for {
			resp, err := r.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				wr.CloseWithError(err)
				return
			}

			wr.Write(resp.Content)
		}
	}()
	err = fileutil.SaveLocalFile(rd, fineName)
	if err != nil {
		return err
	}
	return nil
}

func (c keeperClient) GetBinaryFileList(ctx context.Context) ([]dto.BinaryFile, error) {
	ctx, err := getCtx(c.token)
	if err != nil {
		return nil, err
	}
	r, err := c.client.GetBinaryFileList(ctx, &pb.BinaryFileListRequest{})
	if err != nil {
		return nil, err
	}
	result := make([]dto.BinaryFile, 0, len(r.BinaryFiles))
	for _, f := range r.BinaryFiles {
		result = append(result, dto.BinaryFile{FileName: f.Name, Description: f.Description, Id: uint(f.Id)})
	}
	return result, nil
}
