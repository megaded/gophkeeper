package server

import (
	"context"
	"gophkeeper/internal/config"
	"gophkeeper/internal/identity"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/server/dto"
	"gophkeeper/internal/storage/model"
	pb "gophkeeper/proto"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	storage          Storager
	cfg              config.Config
	userManager      UserManager
	identityProvider identity.IdentityProvider
	pb.UnimplementedKeeperServerServer
}

// AddCredentials implements keeper.KeeperServerServer.
func (s *Server) AddCredentials(context.Context, *pb.AddCredentialsRequest) (*pb.AddCredentialsResponse, error) {
	panic("unimplemented")
}

// Login implements keeper.KeeperServerServer.
func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
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
	return &pb.LoginResponse{Token: token}, nil
}

func (s *Server) Start(ctx context.Context) {
	listen, err := net.Listen("tcp", s.cfg.Address)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()

	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()
	pb.RegisterKeeperServerServer(server, s)
	if err := server.Serve(listen); err != nil {
		panic(err)
	}
}

// Registration implements keeper.KeeperServerServer.
func (s *Server) Registration(ctx context.Context, req *pb.NewUserRequest) (*pb.NewUserResponse, error) {
	resp := &pb.NewUserResponse{}
	return resp, s.userManager.CreateUser(ctx, req.Login, req.Password)
}

var _ pb.KeeperServerServer = (*Server)(nil)

func NewServer(cfg config.Config, storage Storager, userManager UserManager, identityProvider identity.IdentityProvider) Server {
	return Server{storage: storage, cfg: cfg, userManager: userManager, identityProvider: identityProvider}
}

type Storager interface {
	AddUser(ctx context.Context, login string, password string) error
	GetUser(ctx context.Context, login string) (model.User, error)
	AddCredentials(ctx context.Context, cred dto.Credentials) error
	GetCredentials(ctx context.Context, login string) error
	DeleteCredentials(ctx context.Context, login string) error
	UpdateCredentials(ctx context.Context, cred dto.Credentials) error
}

type UserManager interface {
	CreateUser(ctx context.Context, login string, password string) error
}
