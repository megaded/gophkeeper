package server

import (
	"context"
	"gophkeeper/internal/server/dto"
	pb "gophkeeper/proto"
	"io"

	"google.golang.org/grpc"
)

// Загружает текстовые данные
func (s Server) UploadText(ctx context.Context, req *pb.UploadTextRequest) (*pb.UploadTextResponse, error) {
	panic("неа")
}

// Удаляет текстовые данные по Id
func (s Server) DeleteText(ctx context.Context, req *pb.DeleteTextRequest) (*pb.DeleteTextResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	err = s.textManager.DeleteText(ctx, userId, uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteTextResponse{}, nil
}

// Возвращает список текстовый данных
func (s Server) GetTextList(ctx context.Context, req *pb.TextListRequest) (*pb.TextListResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	model, err := s.textManager.GetTextList(ctx, userId)
	if err != nil {
		return nil, err
	}
	data := make([]*pb.Text, 0, len(model))
	for _, f := range data {
		data = append(data, &pb.Text{Id: f.Id, Content: f.Content, IsFile: f.IsFile, Description: f.Description})
	}
	resp := pb.TextListResponse{}
	resp.Info = data
	return &resp, nil
}

func (s Server) UpdateText(ctx context.Context, req *pb.UpdateTextRequest) (*pb.UpdateTextResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	err = s.textManager.UpdateText(ctx, userId, dto.Text{Content: req.Text.Content, Description: req.Text.Description, Id: uint(req.Text.Id)})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateTextResponse{}, nil
}

func (s Server) UpdateTextFile(stream grpc.ClientStreamingServer[pb.UpdateTextFileRequest, pb.UpdateTextFileResponse]) error {
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
	fileInfo, err := s.textManager.GetTextInfo(ctx, uint(req.Id))
	if fileInfo.UserId != userId {
		return err
	}
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

	newFileInfo, err := s.binaryManager.UploadFile(context.Background(), userId, dto.BinaryFile{FileName: req.Filename, Description: req.Description}, rd)
	if err != nil {
		return err
	}
	err = s.textManager.UpdateText(ctx, userId, dto.Text{Id: uint(req.Id), Description: req.Description, IsFile: true, BinaryId: newFileInfo.Id})
	if err != nil {
		err = s.binaryManager.DeleteBinaryFile(ctx, userId, newFileInfo.Id)
		return err
	}

	return stream.SendAndClose(&pb.UpdateTextFileResponse{})
}

func (s Server) UploadTextFile(stream grpc.ClientStreamingServer[pb.UploadTextFileRequest, pb.UploadTextFileResponse]) error {
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

	newFileInfo, err := s.binaryManager.UploadFile(context.Background(), userId, dto.BinaryFile{FileName: req.Filename, Description: req.Description}, rd)
	if err != nil {
		return err
	}
	err = s.textManager.UploadText(ctx, dto.Text{UserId: userId, IsFile: true, Description: req.Description, BinaryId: newFileInfo.Id})
	if err != nil {
		err = s.binaryManager.DeleteBinaryFile(ctx, userId, newFileInfo.Id)
		if err != nil {
			return err
		}
	}
	return stream.SendAndClose(&pb.UploadTextFileResponse{})

}
