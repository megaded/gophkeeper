package proto

import (
	"gophkeeper/internal/config"
	"gophkeeper/internal/logger"
	pb "gophkeeper/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type keeperClient struct {
	client pb.KeeperClient
	token  string
}

func NewKeeperClient() *keeperClient {
	cfg := config.GetConfig()
	conn, err := grpc.NewClient(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	return &keeperClient{client: pb.NewKeeperClient(conn)}
}
