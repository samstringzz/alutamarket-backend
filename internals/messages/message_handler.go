package messages

import (
	"context"

	"github.com/lib/pq"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) CreateChat(ctx context.Context, chatID uint32, usersID pq.Int64Array) error {
	err := h.Service.CreateChat(ctx, chatID, usersID)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) SendMessage(ctx context.Context, req *Message) error {
	err := h.Service.SendMessage(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) GetChatLists(ctx context.Context, userID int64) ([]*Chat, error) {
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
