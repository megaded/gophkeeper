package proto

import (
	"bufio"
	"context"
	fileutil "gophkeeper/internal/client/proto/fileUtil"
	"gophkeeper/internal/dto"
	pb "gophkeeper/proto"
	"io"
	"os"
)

const (
	defaultBufferSize = 32 * 1024   // 32KB
	maxBufferSize     = 1024 * 1024 // 1MB
)

func (c *keeperClient) UploadBinaryFile(ctx context.Context, filePath string, description string) error {
	ctx, err := getCtx(ctx, c.token)
	if err != nil {
		return err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	fStat, err := file.Stat()
	if err != nil {
		return err
	}
	fileName := fStat.Name()
	k, err := c.client.UploadBinaryFile(ctx)
	if err != nil {
		return err
	}

	buf := bufio.NewReader(file)
	data := make([]byte, buf.Size())
	for err != io.EOF {
		_, err := buf.Read(data)
		if err == io.EOF {
			break
		}
		err = k.Send(&pb.UploadBinaryFileRequest{Content: data, Filename: fileName, Description: description})
		if err != nil && err != io.EOF {
			_, err = k.CloseAndRecv()
			return err
		}
	}
	_, err = k.CloseAndRecv()
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func (c *keeperClient) DownloadBinaryFile(ctx context.Context, id uint) error {
	ctx, err := getCtx(ctx, c.token)
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
	defer rd.Close()
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
	ctx, err := getCtx(ctx, c.token)
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

func (c keeperClient) DeleteBinaryFile(ctx context.Context, id uint) error {
	ctx, err := getCtx(ctx, c.token)
	if err != nil {
		return err
	}
	_, err = c.client.DeleteBinaryFile(ctx, &pb.DeleteBinaryFileRequest{Id: uint32(id)})
	return err
}

func (c keeperClient) UpdateBinaryFile(ctx context.Context, filePath string, description string, id uint) error {
	ctx, err := getCtx(ctx, c.token)
	if err != nil {
		return err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	fStat, err := file.Stat()
	if err != nil {
		return err
	}
	fileName := fStat.Name()
	k, err := c.client.UpdateBinaryFile(ctx)
	if err != nil {
		return err
	}

	buf := bufio.NewReader(file)
	bufferSize := defaultBufferSize
	if fStat.Size() > int64(maxBufferSize) {
		bufferSize = maxBufferSize
	}
	data := make([]byte, bufferSize)
	var totalSize int64 = 0
	for err != io.EOF {
		b, err := buf.Read(data)
		totalSize = totalSize + int64(b)
		if err == io.EOF {
			break
		}
		err = k.Send(&pb.UpdateBinaryFileRequest{Content: data, Filename: fileName, Description: description, Id: uint32(id)})
		if err != nil && err != io.EOF {
			_, err = k.CloseAndRecv()
			return err
		}
	}
	_, err = k.CloseAndRecv()
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}
