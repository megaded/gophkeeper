package server

import (
	"context"
	"gophkeeper/internal/config"
	"gophkeeper/internal/identity"
	"gophkeeper/internal/manager"

	"gophkeeper/internal/dto"
	"gophkeeper/internal/server/interceptor"
	"gophkeeper/internal/storage/fileStorage/minio"
	"gophkeeper/internal/storage/model"
	"gophkeeper/internal/storage/postgre"
	pb "gophkeeper/proto"
	"io"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	Storage           Storager
	Cfg               config.Config
	UserManager       UserManager
	IdentityProvider  identity.IdentityProvider
	BinaryManager     BinaryManager
	CreditCardManager CreditCardManager
	CredManager       CredentialsManager
	TextManager       TextManager
	pb.UnimplementedKeeperServer
}

type FileStorager interface {
	UploadFile(ctx context.Context, userName string, fileName string, reader io.Reader, size int64) error
}

func (s *Server) Start(ctx context.Context) {
	listen, err := net.Listen("tcp", s.Cfg.Address)
	if err != nil {
		panic(err)
	}
	authInterceptor := interceptor.GetAuthInterceptor(s.IdentityProvider)
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

func NewServer() (Server, error) {
	cfg := config.GetConfig()
	minio, err := minio.NewStorage(cfg)
	if err != nil {
		return Server{}, err
	}
	storage, err := postgre.NewStorage(&cfg)
	if err != nil {
		return Server{}, err
	}
	builder := NewBuilderServer()
	binaryManager := manager.NewBinaryManager(minio, storage)
	identity := identity.CreateIdentityProvider(&cfg)
	crypto := manager.NewCryptoManager(cfg)
	credManager := manager.NewCredentialsManager(&crypto, storage)
	creditCard := manager.NewCreditCardManager(cfg, storage)
	textManager := manager.NewTextManager(storage)
	builder.SetConfig(cfg)
	builder.SetStorage(storage)
	builder.SetBinaryManager(binaryManager)
	builder.SetIdentityProvider(identity)
	builder.SetCredentialsManager(credManager)
	builder.SetCreditCardManager(creditCard)
	builder.SetTextManager(textManager)
	return builder.Build(), nil
}

type Storager interface {
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

type UserManager interface {
	CreateUser(ctx context.Context, login string, password string) error
}

type BinaryManager interface {
	UploadFile(ctx context.Context, userId uint, dto dto.BinaryFile, reader io.Reader) (dto.BinaryFile, error)
	DownloadFile(ctx context.Context, userId uint, id uint) (reader io.Reader, info dto.BinaryFile, err error)
	GetBinaryFiles(ctx context.Context, userId uint) ([]dto.BinaryFile, error)
	UpdateBinaryFile(ctx context.Context, userId uint, dto dto.BinaryFile, reader io.Reader) error
	DeleteBinaryFile(ctx context.Context, userId uint, id uint) error
}

type CreditCardManager interface {
	UpdateCreditCard(ctx context.Context, userId uint, dto dto.Card) error
	GetCreditCards(ctx context.Context, userId uint) ([]dto.Card, error)
	AddCreditCard(ctx context.Context, userId uint, dto dto.Card) error
	DeleteCreditCard(ctx context.Context, userId uint, id uint) error
}

type CredentialsManager interface {
	AddCredentials(ctx context.Context, userId uint, cred dto.Credentials) error
	GetCredentials(ctx context.Context, userId uint) ([]dto.Credentials, error)
	DeleteCredential(ctx context.Context, userId uint, id uint) error
	UpdateCredentials(ctx context.Context, userId uint, dto dto.Credentials) error
}
type TextManager interface {
	UploadText(ctx context.Context, dto dto.Text) error
	GetTextList(ctx context.Context, userId uint) ([]dto.Text, error)
	GetTextInfo(ctx context.Context, id uint) (dto.Text, error)
	UpdateText(ctx context.Context, userId uint, dto dto.Text) error
	DeleteText(ctx context.Context, userId uint, id uint) error
}
