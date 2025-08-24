package proto

import (
	"context"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/server/dto"
	pb "gophkeeper/proto"
)

func (c *keeperClient) AddCreditCard(ctx context.Context, dto dto.Card) error {
	ctx, err := getCtx(c.token)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	req := &pb.AddCreditCardRequest{Number: dto.Number, Exp: dto.Exp, Cvv: dto.CVV, Description: dto.Description}
	_, err = c.client.AddCreditCard(ctx, req)
	return err
}

func (c *keeperClient) GetCreditCards(ctx context.Context) ([]dto.Card, error) {
	ctx, err := getCtx(c.token)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}
	req := pb.CreditCardRequest{}
	resp, err := c.client.GetCreditCardList(ctx, &req)
	if err != nil {
		return nil, err
	}
	result := make([]dto.Card, 0, len(resp.CreditCards))
	for _, k := range resp.CreditCards {
		result = append(result, dto.Card{Number: k.Number, Exp: k.Exp, CVV: k.Cvv, Description: k.Description})
	}
	return result, nil
}
