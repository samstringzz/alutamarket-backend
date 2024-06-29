package skynet

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

func (s *service) BuyAirtime(c context.Context, req *Airtime) (*string, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.BuyAirtime(ctx, req)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) BuyData(c context.Context, req *Data) (*string, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.BuyData(ctx, req)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) BuyTVSubscription(c context.Context, req *TVSubscription) (*string, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.BuyTVSubscription(ctx, req)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) VerifySmartCard(c context.Context, serviceId, billersCode string) (*SmartcardVerificationResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.VerifySmartCard(ctx, serviceId, billersCode)
	if err != nil {
		return nil, err
	}

	return r, nil
}
