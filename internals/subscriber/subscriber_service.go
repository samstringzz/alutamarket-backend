package subscriber

import (
    "time"
)

type service struct {
    Repository
    timeout time.Duration
}

func NewService(repository Repository) Service {
    return &service{
        Repository: repository,
        timeout:   time.Duration(5) * time.Second,
    }
}

func (s *service) CreateSubscriber(email string) (*Subscriber, error) {
    return s.Repository.CreateSubscriber(email)
}

func (s *service) GetSubscribers() ([]*Subscriber, error) {
    return s.Repository.GetSubscribers()
}