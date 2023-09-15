package cart

import (
	"context"
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

func (s *service) AddToCart(c context.Context, req []*CartItems, user uint32) (*Cart, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	var u []*CartItems
	defer cancel()

	for i, item:= range req{
		u[i] = &CartItems{
		// Product:  req.Product,
		Quantity: item.Quantity,
	}
	}
	r, err := s.Repository.AddToCart(ctx, u, user)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *service) GetCart(c context.Context, user uint32) (*Cart, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	r, err := s.Repository.GetCart(ctx, user)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *service) RemoveFromCart(c context.Context, id uint32) (*Cart, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	r, err := s.Repository.RemoveFromCart(ctx, id)
	if err != nil {
		return nil, err
	}
	return r, nil
}
