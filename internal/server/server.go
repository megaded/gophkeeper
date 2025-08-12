package server

import (
	"context"
	"gophkeeper/internal/config"
	"gophkeeper/internal/server/dto"
)

type Server struct {
	storage    Storager
	cfg        config.Config
	userManger UserManager
}

func NewServer(cfg config.Config, storage Storager, userManager UserManager) Server {
	return Server{storage: storage, cfg: cfg, userManger: userManager}
}

type Storager interface {
	AddUser(ctx context.Context, login string, password string) error
	AddCredentials(ctx context.Context, login string, password string) error
	GetCredentials(ctx context.Context, login string) error
	DeleteCredentials(ctx context.Context, login string) error
	UpdateCredentials(ctx context.Context, cred dto.Credentials)
}

type UserManager interface {
	CreateUser(ctx context.Context, login string, password string) error
}
