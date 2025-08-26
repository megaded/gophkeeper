package proto

import (
	"context"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/server/dto"
	pb "gophkeeper/proto"
)

func (c *keeperClient) AddCreditCard(ctx context.Context, dto dto.Card) error {
	ctx, err := getCtx(ctx, c.token)
	if err != nil {
		return err
	}
	req := &pb.AddCreditCardRequest{Number: dto.Number, Exp: dto.Exp, Cvv: dto.CVV, Description: dto.Description}
	_, err = c.client.AddCreditCard(ctx, req)
	return err
}

func (c *keeperClient) GetCreditCards(ctx context.Context) ([]dto.Card, error) {
	ctx, err := getCtx(ctx, c.token)
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

func (c *keeperClient) DeleteCreditCard(ctx context.Context, id uint) error {
	ctx, err := getCtx(ctx, c.token)
	if err != nil {
		return err
	}
	_, err = c.client.DeleteCreditCard(ctx, &pb.DeleteCreditCardRequest{Id: uint32(id)})
	if err != nil {
		return err
	}
	return nil
}

func (c *keeperClient) UpdateCreditCard(ctx context.Context, dto dto.Card) error {
	ctx, err := getCtx(ctx, c.token)
	if err != nil {
		return err
	}
	req := pb.UpdateCreditCardRequest{
		Card: &pb.CreditCard{Id: uint32(dto.Id), Description: dto.Description, Number: dto.Number, Exp: dto.Exp, Cvv: dto.CVV},
	}
	_, err = c.client.UpdateCreditCard(ctx, &req)
	if err != nil {
		return err
	}
	return nil
}
