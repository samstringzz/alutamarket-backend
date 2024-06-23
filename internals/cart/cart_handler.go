package cart

import (
	"context"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) ModifyCart(ctx context.Context, input *CartItems, user uint32) (*Cart, error) {
	item, err := h.Service.ModifyCart(ctx, input, user)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) GetCart(ctx context.Context, user uint32) (*Cart, error) {
	item, err := h.Service.GetCart(ctx, user)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) RemoveAllCart(ctx context.Context, id uint32) error {
	err := h.Service.RemoveAllCart(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) IntitiatePayment(ctx context.Context, input Order) (string, error) {
	item, err := h.Service.InitiatePayment(ctx, input)
	if err != nil {
		return "", err
	}
	return item, nil
}
