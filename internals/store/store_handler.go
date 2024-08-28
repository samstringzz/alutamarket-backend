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

func (h *Handler) GetStoreByName(ctx context.Context, name string) (*Store, error) {
	r, err := h.Service.GetStoreByName(ctx, name)
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
func (h *Handler) CheckStoreName(ctx context.Context, query string) error {
	err := h.Service.CheckStoreName(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) DeleteStore(ctx context.Context, input uint32) error {
	err := h.Service.DeleteStore(ctx, input)
	return err
}

func (h *Handler) CreateOrder(ctx context.Context, input *StoreOrder) (*StoreOrder, error) {
	item, err := h.Service.CreateOrder(ctx, input)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) UpdateOrder(ctx context.Context, input *StoreOrder) (*StoreOrder, error) {
	item, err := h.Service.UpdateOrder(ctx, input)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) GetOrders(ctx context.Context, storeID uint32) ([]*StoreOrder, error) {
	item, err := h.Service.GetOrders(ctx, storeID)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) GetPurchasedOrders(ctx context.Context, userId string) ([]*Order, error) {
	item, err := h.Service.GetPurchasedOrders(ctx, userId)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) UpdateStoreFollowership(ctx context.Context, storeID uint32, follower Follower, action string) (*Store, error) {
	item, err := h.Service.UpdateStoreFollowership(ctx, storeID, follower, action)
	if err != nil {
		return nil, err
	}
	return item, nil
}
