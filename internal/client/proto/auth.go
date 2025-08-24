package proto

import (
	"context"
	"errors"
	"gophkeeper/internal/logger"
	pb "gophkeeper/proto"

	"google.golang.org/grpc/metadata"
)

func (c *keeperClient) Login(ctx context.Context, login string, password string) (token string, err error) {
	r, err := c.client.Login(ctx, &pb.LoginRequest{Login: login, Password: password})
	if err != nil {
		return "", err
	}
	c.token = r.GetToken()
	return c.token, err
}

func (c *keeperClient) Register(ctx context.Context, login string, password string) error {
	req := &pb.NewUserRequest{Login: login, Password: password}
	_, err := c.client.Registration(context.Background(), req)
	if err != nil {
		logger.Log.Info(err.Error())
	}
	return err
}

func getCtx(token string) (context.Context, error) {
	if token == "" {
		return nil, errors.New("Token is empty")
	}
	md := metadata.New(map[string]string{"token": token})
	return metadata.NewOutgoingContext(context.Background(), md), nil
}
