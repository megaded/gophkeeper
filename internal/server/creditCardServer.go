package server

import (
	"context"
	"gophkeeper/internal/server/dto"
	pb "gophkeeper/proto"
)

// Возвращает список банковских карт пользователя
// Пользователь определяется по переданому токену
func (s *Server) GetCreditCardList(ctx context.Context, req *pb.CreditCardRequest) (*pb.CreditCardListResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	cards, err := s.creditCardManager.GetCreditCards(ctx, userId)
	if err != nil {
		return nil, err
	}
	resp := pb.CreditCardListResponse{}
	resp.CreditCards = make([]*pb.CreditCard, 0, len(cards))
	for _, k := range cards {
		resp.CreditCards = append(resp.CreditCards, &pb.CreditCard{Number: k.Number, Description: k.Description, Cvv: k.CVV, Exp: k.Exp})
	}
	return &resp, err
}

// Удаление банковской карты по ID
// Пользователь определяется по переданому токену
func (s *Server) DeleteCreditCard(ctx context.Context, req *pb.DeleteCreditCardRequest) (*pb.DeleteCreditCardResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	err = s.creditCardManager.DeleteCreditCard(ctx, userId, uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteCreditCardResponse{}, nil
}

// Добавляет тип данных банковская карта
// Пользователь определяется по переданому токену
func (s *Server) AddCreditCard(ctx context.Context, req *pb.AddCreditCardRequest) (*pb.AddCreditCardResponse, error) {
	userId, err := getUserId(ctx)
	if err != nil {
		return nil, err
	}
	err = s.creditCardManager.AddCreditCard(ctx, userId, dto.Card{Number: req.Number, CVV: req.Cvv, Exp: req.Exp, Description: req.Description})
	if err != nil {
		return nil, err
	}
	return &pb.AddCreditCardResponse{}, nil
}
