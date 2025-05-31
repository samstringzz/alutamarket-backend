package store

import (
	"context"
	"fmt"
)

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

func (h *Handler) CreateInvoice(ctx context.Context, input *Invoice) (*Invoice, error) {
	item, err := h.Service.CreateInvoice(ctx, input)
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

func (h *Handler) GetInvoices(ctx context.Context, storeId uint32) ([]*Invoice, error) {
	r, err := h.Service.GetInvoices(ctx, storeId)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (h *Handler) UpdateStore(ctx context.Context, input *UpdateStore) (*Store, error) {
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

func (h *Handler) UpdateOrder(ctx context.Context, input *UpdateStoreOrderInput) (*Order, error) {
	item, err := h.Service.UpdateOrder(ctx, input)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) GetOrders(ctx context.Context, storeID uint32) ([]*Order, error) {
	item, err := h.Service.GetOrders(ctx, storeID)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) GetOrdersByStore(ctx context.Context, storeName string) ([]*Order, error) {
	return h.Service.GetOrdersByStore(ctx, storeName)
}

func (h *Handler) GetPurchasedOrders(ctx context.Context, userId string) ([]*Order, error) {
	item, err := h.Service.GetPurchasedOrders(ctx, userId)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) UpdateStoreFollowership(ctx context.Context, storeID uint32, follower *Follower, action string) (*Store, error) {
	item, err := h.Service.UpdateStoreFollowership(ctx, storeID, follower, action)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) CreateTransactions(ctx context.Context, req *Transactions) (*Transactions, error) {
	item, err := h.Service.CreateTransactions(ctx, req)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) WithdrawFund(ctx context.Context, req *Fund) error {
	err := h.Service.WithdrawFund(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) AddReview(ctx context.Context, req *Review) error {
	err := h.Service.AddReview(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) GetReviews(ctx context.Context, filterType string, value interface{}) ([]*Review, error) {
	result, err := h.Service.GetReviews(ctx, filterType, value)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (h *Handler) GetDVAAccount(ctx context.Context, email string) (*DVAAccount, error) {
	account, err := h.Service.GetDVAAccount(ctx, email)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// In store/handler.go
func (h *Handler) GetFollowedStores(ctx context.Context, userID uint32) ([]*Store, error) {
	return h.Service.GetFollowedStores(ctx, userID)
}

func (h *Handler) GetDVABalance(ctx context.Context, id string) (float64, error) {
	balance, err := h.Service.GetDVABalance(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to get DVA balance: %v", err)
	}
	return balance, nil
}

func (h *Handler) GetOrderByUUID(ctx context.Context, uuid string) (*Order, error) {
	return h.Service.GetOrderByUUID(ctx, uuid)
}

func (h *Handler) UpdateProductUnitsSold(ctx context.Context, productID uint32) error {
	return h.Service.UpdateProductUnitsSold(ctx, productID)
}

func (h *Handler) GetAllOrders(ctx context.Context) ([]*Order, error) {
	return h.Service.GetAllOrders(ctx)
}

func (h *Handler) CheckStoreEarningsDiscrepancy(ctx context.Context, storeID uint32) (int, float64, error) {
	return h.Service.CheckStoreEarningsDiscrepancy(ctx, storeID)
}

// CreatePaystackDVAAccount creates and stores a Paystack DVA account for a store
func (h *Handler) CreatePaystackDVAAccount(ctx context.Context, storeID uint32, account *PaystackDVAResponse, email string) error {
	return h.Service.CreatePaystackDVAAccount(ctx, storeID, account, email)
}

// GetPaystackDVAAccount retrieves a store's Paystack DVA account
func (h *Handler) GetPaystackDVAAccount(ctx context.Context, storeID uint32) (*PaystackDVAResponse, error) {
	return h.Service.GetPaystackDVAAccount(ctx, storeID)
}

// SyncExistingPaystackDVAAccounts retrieves all existing Paystack DVA accounts and stores them in our database
func (h *Handler) SyncExistingPaystackDVAAccounts(ctx context.Context) error {
	return h.Service.SyncExistingPaystackDVAAccounts(ctx)
}
