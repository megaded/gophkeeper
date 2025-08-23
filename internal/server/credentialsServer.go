// Методы по работе с типом данных логин\пароль
package server

import (
	"context"
	"gophkeeper/internal/server/dto"
	pb "gophkeeper/proto"
)

// Создает тип данных логин\пароль
func (s *Server) AddCredentials(ctx context.Context, req *pb.AddCredentialsRequest) (*pb.AddCredentialsResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	err = s.credManager.AddCredentials(ctx, userId, dto.Credentials{Login: req.Login, Password: req.Password, Description: req.Description})
	if err != nil {
		return nil, err
	}
	return &pb.AddCredentialsResponse{}, nil
}

// Возвращает список данных логин\пароль пользователя
// Пользователь определяется по переданому токену
func (s *Server) GetCredentialsList(ctx context.Context, req *pb.CredentialListRequest) (*pb.CredentialListResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	r, err := s.credManager.GetCredentials(ctx, userId)
	if err != nil {
		return nil, err
	}
	result := make([]*pb.Credential, 0, len(r))
	resp := pb.CredentialListResponse{}
	for _, cre := range r {
		result = append(result, &pb.Credential{Login: cre.Login, Password: cre.Password, Description: cre.Description})
	}
	resp.Credentials = result
	return &resp, nil
}

// Удаляет тип данных логин\пароль пользователя по переданому ID
// Пользователь определяется по переданому токену
func (s *Server) DeleteCredential(ctx context.Context, req *pb.DeleteCredentialRequest) (*pb.DeleteCredentialResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	err = s.credManager.DeleteCredential(ctx, userId, uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteCredentialResponse{}, nil
}
