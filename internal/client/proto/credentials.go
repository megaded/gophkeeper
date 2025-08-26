package proto

import (
	"context"
	"gophkeeper/internal/server/dto"
	pb "gophkeeper/proto"
)

func (c *keeperClient) AddCredentials(ctx context.Context, cred dto.Credentials) error {
	req := &pb.AddCredentialsRequest{Login: cred.Login, Password: cred.Password, Description: cred.Description}
	ctx, err := getCtx(ctx, c.token)
	if err != nil {
		return err
	}
	_, err = c.client.AddCredentials(ctx, req)
	return err
}

func (c *keeperClient) GetCredentials(ctx context.Context) ([]dto.Credentials, error) {
	ctx, err := getCtx(ctx, c.token)
	if err != nil {
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

func (c *keeperClient) DeleteCreditial(ctx context.Context, id uint) error {
	ctx, err := getCtx(ctx, c.token)
	if err != nil {
		return err
	}
	_, err = c.client.DeleteCredential(ctx, &pb.DeleteCredentialRequest{Id: uint32(id)})
	if err != nil {
		return err
	}
	return nil
}

func (c keeperClient) UpdateCreditial(ctx context.Context, dto dto.Credentials) error {
	ctx, err := getCtx(ctx, c.token)
	if err != nil {
		return err
	}
	req := pb.UpdateCredentialsRequest{
		Credentials: &pb.Credential{Id: uint32(dto.Id), Description: dto.Description, Login: dto.Login, Password: dto.Password},
	}
	_, err = c.client.UpdateCredentials(ctx, &req)
	if err != nil {
		return err
	}
	return nil
}
