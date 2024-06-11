package store

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

func (s *service) CreateStore(c context.Context, req *Store) (*Store, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.CreateStore(ctx, req)
	if err != nil {
		return nil, err
	}

	return r, nil
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

func (s *service) UpdateStore(c context.Context, req *Store) (*Store, error) {
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
