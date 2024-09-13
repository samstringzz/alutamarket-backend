package messages

import (
	"context"
	"time"

	"github.com/Chrisentech/aluta-market-api/internals/user"
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

func (s *service) FindOrCreateChat(c context.Context, users []*user.User) (*Chat, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	resp, err := s.Repository.FindOrCreateChat(ctx, users)
	if err != nil {
		return nil, err
	}

	return resp, nil
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

func (s *service) GetChatLists(c context.Context, req uint32) ([]*Chat, error) {
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
