package proto

import (
	"bufio"
	"context"
	"errors"
	"gophkeeper/internal/config"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/server/dto"
	pb "gophkeeper/proto"
	"io"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type keeperClient struct {
	client pb.KeeperClient
	token  string
}

func (c *keeperClient) Login(ctx context.Context, login string, password string) (token string, err error) {
	r, err := c.client.Login(ctx, &pb.LoginRequest{Login: login, Password: password})
	if err != nil {
		return "", err
	}
	c.token = r.GetToken()
	return c.token, err
}

func NewKeeperClient() *keeperClient {
	cfg := config.GetConfig()
	conn, err := grpc.NewClient(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	return &keeperClient{client: pb.NewKeeperClient(conn)}
}

func (c *keeperClient) AddCredentials(ctx context.Context, cred dto.Credentials) error {
	req := &pb.AddCredentialsRequest{Login: cred.Login, Password: cred.Password, Description: cred.Description}
	_, err := c.client.AddCredentials(ctx, req)
	return err
}

func (c *keeperClient) AddCreditCard(ctx context.Context, dto dto.Card) error {
	ctx, err := getCtx(c.token)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	req := &pb.AddCreditCardRequest{Number: dto.Number, Exp: dto.Exp, Cvv: dto.CVV, Description: dto.Description}
	_, err = c.client.AddCreditCard(ctx, req)
	return err
}

func (c *keeperClient) GetCreditCards(ctx context.Context) ([]dto.Card, error) {
	ctx, err := getCtx(c.token)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}
	req := pb.CreditCardRequest{}
	resp, err := c.client.GetCreditCardList(ctx, &req)
	if err != nil {
		return nil, err
	}
	result := make([]dto.Card, 0, len(resp.CreditCards))
	for _, k := range resp.CreditCards {
		result = append(result, dto.Card{Number: k.Number, Exp: k.Exp, CVV: k.Cvv, Description: k.Description})
	}
	return result, nil
}

func (c *keeperClient) Register(ctx context.Context, login string, password string) error {
	req := &pb.NewUserRequest{Login: login, Password: password}
	_, err := c.client.Registration(context.Background(), req)
	if err != nil {
		logger.Log.Info(err.Error())
	}
	return err
}

func (c *keeperClient) UploadBinaryFile(reader io.Reader, fileName string, description string, size int64) error {
	ctx, err := getCtx(c.token)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	k, err := c.client.UploadBinaryFile(ctx)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	buf := bufio.NewReader(reader)
	data := make([]byte, buf.Size())
	var totalSize int64 = 0
	for {
		b, err := buf.Read(data)
		totalSize = totalSize + int64(b)
		if err == io.EOF {
			logger.Log.Info("байтов", zap.Int64("всего", totalSize))
			logger.Log.Info("Конец файла")
			break
		}
		err = k.Send(&pb.UploadBinaryFileRequest{Content: data, Filename: fileName, Description: description, Size: size})
		if err != nil && err != io.EOF {
			_, err = k.CloseAndRecv()
			logger.Log.Error(err.Error())
			return err
		}
	}
	_, err = k.CloseAndRecv()
	logger.Log.Info("Закончили отправку")
	if err != nil && err != io.EOF {
		logger.Log.Error(err.Error())
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
