// Методы сервера авторизации и аутентификации
package server

import (
	"context"
	"errors"
	"gophkeeper/internal/internal_error"
	"gophkeeper/internal/server/interceptor"
	pb "gophkeeper/proto"
	"strconv"

	"google.golang.org/grpc/metadata"
)

// Аутентификация пользователя по логину и паролю
// Возвращает токен пользователя
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

// Регистрация нового пользователя
func (s *Server) Registration(ctx context.Context, req *pb.NewUserRequest) (*pb.NewUserResponse, error) {
	resp := &pb.NewUserResponse{}
	return resp, s.userManager.CreateUser(ctx, req.Login, req.Password)
}

func getUserId(ctx context.Context) (uint, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	var userId string
	if ok {
		values := md.Get(interceptor.UserId)
		if len(values) > 0 {
			// ключ содержит слайс строк, получаем первую строку
			userId = values[0]
		}
		if userId == "" {
			return 0, errors.New("UserId empty")
		}
	}
	userid, err := strconv.ParseUint(userId, 10, 32)
	return uint(userid), err
}
