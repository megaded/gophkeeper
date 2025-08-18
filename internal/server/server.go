package server

import (
	"context"
	"gophkeeper/internal/config"
	"gophkeeper/internal/identity"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/server/dto"
	"gophkeeper/internal/storage/model"
	keeper "gophkeeper/proto"
)

type Server struct {
	storage          Storager
	cfg              config.Config
	userManger       UserManager
	identityProvider identity.IdentityProvider
	keeper.UnimplementedKeeperServerServer
}

// AddCredentials implements keeper.KeeperServerServer.
func (s *Server) AddCredentials(context.Context, *keeper.AddCredentialsRequest) (*keeper.AddCredentialsResponse, error) {
	panic("unimplemented")
}

// Login implements keeper.KeeperServerServer.
func (s *Server) Login(ctx context.Context, req *keeper.LoginRequest) (*keeper.LoginResponse, error) {
	userInfo, err := s.storage.GetUser(ctx, req.Login)
	if err != nil {
		return nil, err
	}
	ok := s.identityProvider.VerifyPassword(userInfo.Hash, req.Password)
	if !ok {
		return nil, internal_error.ErrInvalidPassword
	}
	token, err := s.identityProvider.GenerateToken(int(userInfo.ID))
	if err != nil {
		return nil, err
	}
	return &keeper.LoginResponse{Token: token}, nil
}

// Registration implements keeper.KeeperServerServer.
func (s *Server) Registration(ctx context.Context, req *keeper.NewUserRequest) (*keeper.NewUserResponse, error) {
	return &keeper.NewUserResponse{}, s.userManger.CreateUser(ctx, req.Login, req.Password)
}

var _ keeper.KeeperServerServer = (*Server)(nil)

func NewServer(cfg config.Config, storage Storager, userManager UserManager, identityProvider identity.IdentityProvider) Server {
	return Server{storage: storage, cfg: cfg, userManger: userManager, identityProvider: identityProvider}
}

type Storager interface {
	AddUser(ctx context.Context, login string, password string) error
	GetUser(ctx context.Context, login string) (model.User, error)
	AddCredentials(ctx context.Context, login string, password string) error
	GetCredentials(ctx context.Context, login string) error
	DeleteCredentials(ctx context.Context, login string) error
	UpdateCredentials(ctx context.Context, cred dto.Credentials)
}

type UserManager interface {
	CreateUser(ctx context.Context, login string, password string) error
}
