package proto

import (
	"context"
	pb "gophkeeper/proto"
)

type keeperClient struct {
	client pb.KeeperServerClient
}

func (c *keeperClient) Login(ctx context.Context, login string, password string) (token string, err error) {
	r, err := c.client.Login(ctx, &pb.LoginRequest{Login: login, Password: password})
	if err != nil {
		return "", err
	}
	return r.GetToken(), err
}

func NewKeeperClient() keeperClient {
	return keeperClient{}
}

func (c *keeperClient) Register(ctx context.Context, login string, password string) error {
	_, err := c.client.Registration(ctx, &pb.NewUserRequest{Login: login, Password: password})
	return err
}
