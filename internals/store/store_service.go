package store

import (
	"context"
	"fmt"
	"time"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(5) * time.Second,
	}
}

func (s *service) CreateStore(c context.Context, req *Store) (*Store, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.CreateStore(ctx, req)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) CreateInvoice(c context.Context, req *Invoice) (*Invoice, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.CreateInvoice(ctx, req)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) CheckStoreName(c context.Context, query string) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	err := s.Repository.CheckStoreName(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetStore(c context.Context, id uint32) (*Store, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.GetStore(ctx, id)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) GetStoreByName(c context.Context, name string) (*Store, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.GetStoreByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) GetStores(ctx context.Context, user uint32, limit, offset int) ([]*Store, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	r, err := s.Repository.GetStores(ctx, user, limit, offset)
	if err != nil {
		return nil, err
	}

	return r, nil
}
func (s *service) GetInvoices(ctx context.Context, storeId uint32) ([]*Invoice, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	r, err := s.Repository.GetInvoices(ctx, storeId)
	if err != nil {
		return nil, err
	}

	return r, nil
}
func (s *service) UpdateStore(c context.Context, req *UpdateStore) (*Store, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.UpdateStore(ctx, req)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) DeleteStore(c context.Context, id uint32) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	err := s.Repository.DeleteStore(ctx, id)
	return err
}

func (s *service) CreateOrder(c context.Context, req *StoreOrder) (*StoreOrder, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	resp, err := s.Repository.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *service) GetOrders(c context.Context, storeId uint32) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	resp, err := s.Repository.GetOrders(ctx, storeId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *service) GetPurchasedOrders(c context.Context, userId string) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	resp, err := s.Repository.GetPurchasedOrders(ctx, userId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *service) GetOrdersByStore(ctx context.Context, storeName string) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.Repository.GetOrdersByStore(ctx, storeName)
}

func (s *service) UpdateOrderStatus(ctx context.Context, uuid string, status, transStatus string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.Repository.UpdateOrderStatus(ctx, uuid, status, transStatus)
}

func (s *service) UpdateOrder(c context.Context, req *UpdateStoreOrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	resp, err := s.Repository.UpdateOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *service) UpdateStoreFollowership(ctx context.Context, storeID uint32, follower *Follower, action string) (*Store, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	store, err := s.Repository.UpdateStoreFollowership(ctx, storeID, follower, action)
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (s *service) CreateTransactions(ctx context.Context, req *Transactions) (*Transactions, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	store, err := s.Repository.CreateTransactions(ctx, req)
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (s *service) WithdrawFund(ctx context.Context, req *Fund) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.Repository.WithdrawFund(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AddReview(ctx context.Context, review *Review) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	if err := s.Repository.AddReview(ctx, review); err != nil {
		return err
	}
	return nil
}

func (s *service) GetReviews(ctx context.Context, filterType string, value interface{}) ([]*Review, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	result, err := s.Repository.GetReviews(ctx, filterType, value)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) GetDVAAccount(ctx context.Context, email string) (*DVAAccount, error) {
	return s.Repository.GetDVAAccount(ctx, email)
}

func (s *service) GetDVABalance(ctx context.Context, id string) (float64, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	balance, err := s.Repository.GetDVABalance(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to get DVA balance: %v", err)
	}
	return balance, nil
}

func (s *service) GetFollowedStores(ctx context.Context, userID uint32) ([]*Store, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stores, err := s.Repository.GetFollowedStores(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get followed stores: %v", err)
	}
	return stores, nil
}

func (s *service) GetOrderByUUID(ctx context.Context, uuid string) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.Repository.GetOrderByUUID(ctx, uuid)
}

func (s *service) UpdateProductUnitsSold(ctx context.Context, productID uint32) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.Repository.UpdateProductUnitsSold(ctx, productID)
}
