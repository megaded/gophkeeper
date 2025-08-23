package server

import (
	"context"
	"gophkeeper/internal/logger"
	pb "gophkeeper/proto"
	"io"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Загружает бинарный файл произвольной длинны до 5 гигабайт
// Пользователь определяется по переданому токену
func (s *Server) UploadBinaryFile(stream grpc.ClientStreamingServer[pb.UploadBinaryFileRequest, pb.UploadBinaryFileResponse]) error {
	ctx := stream.Context()
	userId, err := getUserId(ctx)
	if err != nil {
		return err
	}
	rd, wr := io.Pipe()
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	fileName := req.Filename
	size := req.Size
	var totalSize int64 = int64(len(req.Content))
	defer rd.Close()

	go func() {
		wr.Write(req.Content)
		defer wr.Close()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				logger.Log.Info("Пришло", zap.Int64("байт", totalSize))
				logger.Log.Info("Выходим")
				return
			}
			if err != nil {
				wr.CloseWithError(err)
				return
			}
			totalSize = totalSize + int64(len(req.Content))

			wr.Write(req.Content)
		}
	}()

	err = s.fileStorage.UploadFile(context.Background(), string(userId), fileName, rd, size)
	if err != nil {
		return err
	}
	logger.Log.Info("Загрузили что-то")
	return stream.SendAndClose(&pb.UploadBinaryFileResponse{})
}

// Загружает текстовый файл произвольной длинны до 5 гигабайт
// Пользователь определяется по переданому токену
func (s Server) UploadTextFile(grpc.ClientStreamingServer[pb.UploadTextFileRequest, pb.UploadTextFileRequest]) error {
	panic("неаы")
}
