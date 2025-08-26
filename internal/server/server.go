package server

import (
	"context"
	"gophkeeper/internal/config"
	"gophkeeper/internal/identity"
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
	userManager       userManager
	identityProvider  identity.IdentityProvider
	binaryManager     binaryManager
	creditCardManager creditCardManager
	credManager       credentialsManager
	textManager       textManager
	pb.UnimplementedKeeperServer
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
	server := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.UnaryAuthInterceptor), grpc.StreamInterceptor(authInterceptor.StreamingAuthInterceptor))

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

func NewServer(cfg config.Config, storage storager, userManager userManager,
	identityProvider identity.IdentityProvider, binaryManager binaryManager,
	creditCardManager creditCardManager, credManager credentialsManager,
	textManager textManager) Server {
	return Server{storage: storage,
		cfg:               cfg,
		userManager:       userManager,
		identityProvider:  identityProvider,
		binaryManager:     binaryManager,
		creditCardManager: creditCardManager,
		credManager:       credManager,
		textManager:       textManager}
}

type storager interface {
	userStorager
	credentialsStorager
	creditCardStorager
	textStorager
}

type textStorager interface {
	GetText(ctx context.Context, id uint) (model.Text, error)
	GetTextList(ctx context.Context, userId uint) ([]model.Text, error)
	UpdateText(ctx context.Context, id uint, content string, description string) error
	DeleteText(ctx context.Context, id uint) error
}

type credentialsStorager interface {
	AddCredentials(ctx context.Context, userId uint, login []byte, password []byte, description string) error
	GetCredentials(ctx context.Context, userId uint) ([]model.Credentials, error)
	DeleteCredentials(ctx context.Context, id uint) error
	UpdateCredentials(ctx context.Context, id uint, login []byte, password []byte, description string) error
}

type userStorager interface {
	AddUser(ctx context.Context, login string, password string) error
	GetUser(ctx context.Context, login string) (model.User, error)
}

type creditCardStorager interface {
	AddCreditCard(ctx context.Context, userId uint, number []byte, ext []byte, cvv []byte, description string) error
	DeleteCreditCard(ctx context.Context, id uint) error
	GetCreditCards(ctx context.Context, userId uint) ([]model.CreditCard, error)
	UpdateCreditCard(ctx context.Context, id uint, cvv []byte, exp []byte, cve []byte, description string) error
}

type userManager interface {
	CreateUser(ctx context.Context, login string, password string) error
}

type binaryManager interface {
	UploadFile(ctx context.Context, userId uint, dto dto.BinaryFile, reader io.Reader) (dto.BinaryFile, error)
	DownloadFile(ctx context.Context, userId uint, id uint) (reader io.Reader, info dto.BinaryFile, err error)
	GetBinaryFiles(ctx context.Context, userId uint) ([]dto.BinaryFile, error)
	UpdateBinaryFile(ctx context.Context, userId uint, dto dto.BinaryFile, reader io.Reader) error
	DeleteBinaryFile(ctx context.Context, userId uint, id uint) error
}

type creditCardManager interface {
	UpdateCreditCard(ctx context.Context, userId uint, dto dto.Card) error
	GetCreditCards(ctx context.Context, userId uint) ([]dto.Card, error)
	AddCreditCard(ctx context.Context, userId uint, dto dto.Card) error
	DeleteCreditCard(ctx context.Context, userId uint, id uint) error
}

type credentialsManager interface {
	AddCredentials(ctx context.Context, userId uint, cred dto.Credentials) error
	GetCredentials(ctx context.Context, userId uint) ([]dto.Credentials, error)
	DeleteCredential(ctx context.Context, userId uint, id uint) error
	UpdateCredentials(ctx context.Context, userId uint, dto dto.Credentials) error
}
type textManager interface {
	UploadText(ctx context.Context, dto dto.Text) error
	GetTextList(ctx context.Context, userId uint) ([]dto.Text, error)
	GetTextInfo(ctx context.Context, id uint) (dto.Text, error)
	UpdateText(ctx context.Context, userId uint, dto dto.Text) error
	DeleteText(ctx context.Context, userId uint, id uint) error
}
