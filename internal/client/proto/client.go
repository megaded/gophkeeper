package proto

import (
	"context"
	"gophkeeper/internal/logger"
	pb "gophkeeper/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type keeperClient struct {
	client pb.KeeperServerClient
}

func (c keeperClient) Login(ctx context.Context, login string, password string) (token string, err error) {
	r, err := c.client.Login(ctx, &pb.LoginRequest{Login: login, Password: password})
	if err != nil {
		return "", err
	}
	return r.GetToken(), err
}

func NewKeeperClient() *keeperClient {
	//cfg := config.GetConfig()
	conn, err := grpc.NewClient(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	return &keeperClient{client: pb.NewKeeperServerClient(conn)}
}

func (c keeperClient) Register(ctx context.Context, login string, password string) error {
	req := &pb.NewUserRequest{Login: login, Password: password}
	_, err := c.client.Registration(context.Background(), req)
	if err != nil {
		logger.Log.Info(err.Error())
	}
	return err
}
