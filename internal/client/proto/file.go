package proto

import (
	"context"
	"gophkeeper/internal/server/dto"
	pb "gophkeeper/proto"
)

func (c *keeperClient) UploadTextFile(ctx context.Context, fileName string, description string) error {
	return nil
}

func (c *keeperClient) UpdateTextFile(ctx context.Context, filePath string, description string, id uint) error {
	ctx, err := getCtx(ctx, c.token)
	if err != nil {
		return err
	}
	return nil
}
func (c *keeperClient) UploadText(ctx context.Context, content string, description string) error {
	ctx, err := getCtx(ctx, c.token)
	if err != nil {
		return err
	}
	req := pb.UploadTextRequest{Content: content, Description: description}
	_, err = c.client.UploadText(ctx, &req)
	if err != nil {
		return err
	}

	return nil
}
func (c *keeperClient) GetTextList(ctx context.Context) ([]dto.Text, error) {
	ctx, err := getCtx(ctx, c.token)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.GetTextList(ctx, &pb.TextListRequest{})
	if err != nil {
		return nil, err
	}
	data := make([]dto.Text, 0, len(resp.Info))
	for _, k := range resp.Info {
		data = append(data, dto.Text{Content: k.Content, Id: uint(k.Id), Description: k.Description, IsFile: k.IsFile})
	}
	return data, nil
}
func (c *keeperClient) DeleteText(ctx context.Context, id uint) error {
	ctx, err := getCtx(ctx, c.token)
	if err != nil {
		return err
	}
	req := pb.DeleteTextRequest{Id: uint32(id)}
	_, err = c.client.DeleteText(ctx, &req)
	if err != nil {
		return err
	}
	return nil
}
func (c *keeperClient) UpdateText(ctx context.Context, dto dto.Text) error {
	ctx, err := getCtx(ctx, c.token)
	if err != nil {
		return err
	}
	req := pb.UpdateTextRequest{Text: &pb.Text{Description: dto.Description, Content: dto.Content, Id: uint32(dto.Id)}}
	_, err = c.client.UpdateText(ctx, &req)
	if err != nil {
		return err
	}
	return nil
}
