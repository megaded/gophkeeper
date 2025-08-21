package server

import (
	"context"
	"gophkeeper/internal/config"
	"gophkeeper/internal/identity"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/server/dto"
	"gophkeeper/internal/server/interceptor"
	"gophkeeper/internal/storage/model"
	pb "gophkeeper/proto"
	"io"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	storage          Storager
	cfg              config.Config
	userManager      UserManager
	identityProvider identity.IdentityProvider
	fileStorage      FileStorager
	pb.UnimplementedKeeperServer
}

// DownloadBinaryFile implements keeper.KeeperServer.
func (s *Server) DownloadBinaryFile(context.Context, *pb.DownloadBinaryFileRequest) (*pb.UploadBinaryFileResponse, error) {
	panic("unimplemented")
}

type FileStorager interface {
	UploadFile(ctx context.Context, userName string, fileName string, reader io.Reader, size int64) error
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

func (s *Server) UploadBinaryFile(stream grpc.ClientStreamingServer[pb.UploadBinaryFileRequest, pb.UploadBinaryFileResponse]) error {
	ctx := stream.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	var userId string
	if ok {
		values := md.Get(interceptor.UserId)
		if len(values) > 0 {
			// ключ содержит слайс строк, получаем первую строку
			userId = values[0]
		}
		userId = "keeperrr"
		if userId == "" {
			return nil
		}
	}

	rd, wr := io.Pipe()
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	fileName := req.Filename
	size := req.Size
	var totalSize int64 = int64(len(req.Content))
	defer rd.Close()

	go func() {
		wr.Write(req.Content)
		defer wr.Close()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				logger.Log.Info("Пришло", zap.Int64("байт", totalSize))
				logger.Log.Info("Выходим")
				return
			}
			if err != nil {
				wr.CloseWithError(err)
				return
			}
			totalSize = totalSize + int64(len(req.Content))

			wr.Write(req.Content)
		}
	}()

	err = s.fileStorage.UploadFile(context.Background(), userId, fileName, rd, size)
	if err != nil {
		return err
	}
	logger.Log.Info("Загрузили что-то")
	return stream.SendAndClose(&pb.UploadBinaryFileResponse{})
}

func (s *Server) Start(ctx context.Context) {
	listen, err := net.Listen("tcp", s.cfg.Address)
	if err != nil {
		panic(err)
	}
	//authInterceptor := interceptor.GetAuthInterceptor(s.identityProvider)
	server := grpc.NewServer()

	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()
	pb.RegisterKeeperServer(server, s)
	if err := server.Serve(listen); err != nil {
		panic(err)
	}
}

// Registration implements keeper.KeeperServerServer.
func (s *Server) Registration(ctx context.Context, req *pb.NewUserRequest) (*pb.NewUserResponse, error) {
	resp := &pb.NewUserResponse{}
	return resp, s.userManager.CreateUser(ctx, req.Login, req.Password)
}

var _ pb.KeeperServer = (*Server)(nil)

func NewServer(cfg config.Config, storage Storager, userManager UserManager, identityProvider identity.IdentityProvider, fileStorage FileStorager) Server {
	return Server{storage: storage, cfg: cfg, userManager: userManager, identityProvider: identityProvider, fileStorage: fileStorage}
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
