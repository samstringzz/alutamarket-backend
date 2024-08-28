package cart

import (
	"context"
	"errors"
	"time"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(20) * time.Second,
	}
}

func (s *service) ModifyCart(c context.Context, req *CartItems, user uint32) (*Cart, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.ModifyCart(ctx, req, user)
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

func (s *service) RemoveAllCart(c context.Context, id uint32) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	err := s.Repository.RemoveAllCart(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) InitiatePayment(c context.Context, input Order) (string, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	r, err := s.Repository.InitiatePayment(ctx, input)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", errors.New("operation timed out: payment initiation took too long")
		}
		return "", err
	}
	return r, nil
}
