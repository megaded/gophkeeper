package proto

import (
	"context"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/server/dto"
	pb "gophkeeper/proto"
)

func (c *keeperClient) AddCredentials(ctx context.Context, cred dto.Credentials) error {
	req := &pb.AddCredentialsRequest{Login: cred.Login, Password: cred.Password, Description: cred.Description}
	_, err := c.client.AddCredentials(ctx, req)
	return err
}

func (c *keeperClient) GetCredentials(ctx context.Context) ([]dto.Credentials, error) {
	ctx, err := getCtx(c.token)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}
	req := pb.CredentialListRequest{}
	resp, err := c.client.GetCredentialsList(ctx, &req)
	if err != nil {
		return nil, err
	}
	result := make([]dto.Credentials, 0, len(resp.Credentials))
	for _, k := range resp.Credentials {
		result = append(result, dto.Credentials{Login: k.Login, Password: k.Password, Description: k.Description})
	}
	return result, nil
}
