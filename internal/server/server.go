package server

import (
	"context"
	"gophkeeper/internal/config"
	"gophkeeper/internal/identity"
	"gophkeeper/internal/manager"
	"gophkeeper/internal/server/dto"
	"gophkeeper/internal/server/interceptor"
	"gophkeeper/internal/storage/model"
	pb "gophkeeper/proto"
	"io"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	storage           storager
	cfg               config.Config
	userManager       UserManager
	identityProvider  identity.IdentityProvider
	fileStorage       FileStorager
	creditCardManager manager.CreditCardManager
	credManager       manager.CredentialsManager
	pb.UnimplementedKeeperServer
}

func (s *Server) DownloadBinaryFile(context.Context, *pb.DownloadBinaryFileRequest) (*pb.UploadBinaryFileResponse, error) {
	panic("unimplemented")
}

type FileStorager interface {
	UploadFile(ctx context.Context, userName string, fileName string, reader io.Reader, size int64) error
}

func (s *Server) Start(ctx context.Context) {
	listen, err := net.Listen("tcp", s.cfg.Address)
	if err != nil {
		panic(err)
	}
	authInterceptor := interceptor.GetAuthInterceptor(s.identityProvider)
	server := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.UnaryInterceptor))

	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()
	pb.RegisterKeeperServer(server, s)
	if err := server.Serve(listen); err != nil {
		panic(err)
	}
}

var _ pb.KeeperServer = (*Server)(nil)

func NewServer(cfg config.Config, storage storager, userManager UserManager,
	identityProvider identity.IdentityProvider, fileStorage FileStorager,
	creditCardManager manager.CreditCardManager, credManager manager.CredentialsManager) Server {
	return Server{storage: storage,
		cfg:               cfg,
		userManager:       userManager,
		identityProvider:  identityProvider,
		fileStorage:       fileStorage,
		creditCardManager: creditCardManager,
		credManager:       credManager}
}

type storager interface {
	userStorager
	credentialsStorager
	creditCardStorager
}

type fileStorager interface {
	AddText(ctx context.Context, userId uint, content string, description string) error
}

type credentialsStorager interface {
	AddCredentials(ctx context.Context, userId uint, login []byte, password []byte) error
	GetCredentials(ctx context.Context, userId uint) ([]model.Credentials, error)
	DeleteCredentials(ctx context.Context, id uint) error
	UpdateCredentials(ctx context.Context, cred dto.Credentials) error
}

type userStorager interface {
	AddUser(ctx context.Context, login string, password string) error
	GetUser(ctx context.Context, login string) (model.User, error)
}

type creditCardStorager interface {
	AddCreditCard(ctx context.Context, userId uint, number []byte, ext []byte, cvv []byte, description string) error
	DeleteCreditCard(ctx context.Context, id uint) error
	GetCreditCards(ctx context.Context, userId uint) ([]model.CreditCard, error)
}

type UserManager interface {
	CreateUser(ctx context.Context, login string, password string) error
}
