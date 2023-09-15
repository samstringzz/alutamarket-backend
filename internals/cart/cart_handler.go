package cart

import "context"

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) AddToCart(ctx context.Context, input []*CartItems, user uint32) (*Cart, error) {
	item, err := h.Service.AddToCart(ctx, input, user)
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

func (h *Handler) RemoveFromCart(ctx context.Context, id uint32) (*Cart, error) {
	item, err := h.Service.RemoveFromCart(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}
