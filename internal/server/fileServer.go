package server

import (
	"context"
	"gophkeeper/internal/dto"
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

	var totalSize int64 = int64(len(req.Content))
	defer rd.Close()

	go func() {
		wr.Write(req.Content)
		defer wr.Close()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
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

	_, err = s.BinaryManager.UploadFile(context.Background(), userId, dto.BinaryFile{FileName: req.Filename, Description: req.Description}, rd)
	if err != nil {
		return err
	}
	return stream.SendAndClose(&pb.UploadBinaryFileResponse{})
}

// Загрузка файла
func (s *Server) DownloadBinaryFile(req *pb.DownloadBinaryFileRequest, resp grpc.ServerStreamingServer[pb.DownloadBinaryFileResponse]) error {
	logger.Log.Info("Начинаем загрузку файла", zap.Int64("id", req.Id))
	ctx := resp.Context()
	userId, err := getUserId(ctx)
	if err != nil {
		return err
	}
	reader, meta, err := s.BinaryManager.DownloadFile(ctx, userId, uint(req.Id))
	if err != nil {
		return err
	}

	bufferSize := int64(10 * 1024)
	buffer := make([]byte, bufferSize)

	for {

		n, err := reader.Read(buffer)
		if n > 0 {

			sendErr := resp.Send(&pb.DownloadBinaryFileResponse{Content: buffer, Filename: meta.FileName})
			if sendErr != nil {
				return sendErr
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}
	}
	return nil
}

// Получение списка файлов
func (s *Server) GetBinaryFileList(ctx context.Context, request *pb.BinaryFileListRequest) (*pb.BinaryFileListResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	r, err := s.BinaryManager.GetBinaryFiles(ctx, userId)
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

func (s Server) DeleteBinaryFile(ctx context.Context, req *pb.DeleteBinaryFileRequest) (*pb.DeleteBinaryFileResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	err = s.BinaryManager.DeleteBinaryFile(ctx, userId, uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteBinaryFileResponse{}, nil
}

func (s Server) UpdateBinaryFile(stream grpc.ClientStreamingServer[pb.UpdateBinaryFileRequest, pb.UpdateBinaryFileResponse]) error {
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

	err = s.BinaryManager.UpdateBinaryFile(context.Background(), userId, dto.BinaryFile{FileName: req.Filename, Description: req.Description, Id: uint(req.Id)}, rd)
	if err != nil {
		return err
	}
	return stream.SendAndClose(&pb.UpdateBinaryFileResponse{})
}
