package store

import "context"

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) CreateStore(ctx context.Context, input *Store) (*Store, error) {
	item, err := h.Service.CreateStore(ctx, input)
	if err != nil {
		return nil, err
	}
	return item, nil
}
func (h *Handler) GetStore(ctx context.Context, input uint32) (*Store, error) {
	r, err := h.Service.GetStore(ctx, input)
	if err != nil {
		return nil, err
	}
	return r, nil
}
func (h *Handler) GetStores(ctx context.Context, user uint32, limit, offset int) ([]*Store, error) {
	r, err := h.Service.GetStores(ctx, user, limit, offset)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (h *Handler) UpdateStore(ctx context.Context, input *Store) (*Store, error) {
	item, err := h.Service.UpdateStore(ctx, input)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) DeleteStore(ctx context.Context, input uint32) error {
	err := h.Service.DeleteStore(ctx, input)
	return err
}
