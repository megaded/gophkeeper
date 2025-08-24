package server

import (
	"bufio"
	"context"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/server/dto"
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

	err = s.binaryManager.UploadFile(context.Background(), userId, dto.BinaryFile{FileName: req.Filename, Description: req.Description}, rd)
	if err != nil {
		return err
	}
	logger.Log.Info("Загрузили что-то")
	return stream.SendAndClose(&pb.UploadBinaryFileResponse{})
}

func (s *Server) DownloadBinaryFile(req *pb.DownloadBinaryFileRequest, resp grpc.ServerStreamingServer[pb.DownloadBinaryFileResponse]) error {
	logger.Log.Info("Начинаем загрузку файла %s", zap.Int64("id", req.Id))
	ctx := resp.Context()
	userId, err := getUserId(ctx)
	if err != nil {
		return err
	}
	reader, meta, err := s.binaryManager.DownloadFile(ctx, userId, uint(req.Id))
	if err != nil {
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
		err = resp.Send(&pb.DownloadBinaryFileResponse{Content: data, Filename: meta.FileName})

	}

	return nil
}

// Загружает текстовый файл произвольной длинны до 5 гигабайт
// Пользователь определяется по переданому токену
func (s Server) UploadTextFile(grpc.ClientStreamingServer[pb.UploadTextFileRequest, pb.UploadTextFileRequest]) error {
	panic("неаы")
}
