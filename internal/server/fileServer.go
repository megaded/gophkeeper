package server

import (
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
	logger.Log.Info("Начинаем загрузку файла", zap.Int64("id", req.Id))
	ctx := resp.Context()
	userId, err := getUserId(ctx)
	if err != nil {
		return err
	}
	reader, meta, err := s.binaryManager.DownloadFile(ctx, userId, uint(req.Id))
	if err != nil {
		return err
	}
	data := make([]byte, 4000)

	for err != io.EOF {
		_, err := reader.Read(data)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			return err
		}
		err = resp.Send(&pb.DownloadBinaryFileResponse{Content: data, Filename: meta.FileName})
		if err != nil {
			return err
		}
	}
	logger.Log.Info("Закончили отправку")
	return nil
}

func (s *Server) GetBinaryFileList(ctx context.Context, request *pb.BinaryFileListRequest) (*pb.BinaryFileListResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	r, err := s.binaryManager.GetBinaryFiles(ctx, userId)
	if err != nil {
		return nil, err
	}
	result := make([]*pb.BinaryFile, 0, len(r))
	resp := pb.BinaryFileListResponse{}
	for _, f := range r {
		result = append(result, &pb.BinaryFile{Id: uint32(f.Id), Name: f.FileName, Description: f.Description})
	}
	resp.BinaryFiles = result
	return &resp, nil
}

// Загружает текстовый файл произвольной длинны до 5 гигабайт
// Пользователь определяется по переданому токену
func (s Server) UploadTextFile(grpc.ClientStreamingServer[pb.UploadTextFileRequest, pb.UploadTextFileRequest]) error {
	panic("неаы")
}
