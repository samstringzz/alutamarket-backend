package messages

import (
	"context"
	"time"

	"github.com/lib/pq"
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

func (s *service) CreateChat(c context.Context, chatID uint32, usersID pq.Int64Array) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	err := s.Repository.CreateChat(ctx, chatID, usersID)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) SendMessage(c context.Context, req *Message) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	err := s.Repository.SendMessage(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetChatLists(c context.Context, req int64) ([]*Chat, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	chats, err := s.Repository.GetChatLists(ctx, req)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (s *service) GetMessages(c context.Context, req string) ([]*Message, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	chats, err := s.Repository.GetMessages(ctx, req)
	if err != nil {
		return nil, err
	}

	return chats, nil
}
