package messages

import (
	"context"

	"github.com/Chrisentech/aluta-market-api/internals/user"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) FindOrCreateChat(ctx context.Context, users []*user.User) (*Chat, error) {
	resp, err := h.Service.FindOrCreateChat(ctx, users)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *Handler) SendMessage(ctx context.Context, req *Message) error {
	err := h.Service.SendMessage(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) GetChatLists(ctx context.Context, userID uint32) ([]*Chat, error) {
	chats, err := h.Service.GetChatLists(ctx, userID)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (h *Handler) GetMessages(ctx context.Context, chatId string) ([]*Message, error) {
	messages, err := h.Service.GetMessages(ctx, chatId)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
