package messages

import (
	"context"
	"log"
	"net/http"

	"github.com/Chrisentech/aluta-market-api/internals/user"
)

type Handler struct {
	Service
}

func NewMessageHandler(s Service) *Handler {
	if s == nil {
		log.Println("MessageService passed to NewMessageHandler is nil!")
		return nil
	}
	handler := &Handler{Service: s}
	log.Printf("Message Handler created successfully: %+v", handler)
	return handler
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

// Add this method to Handler
func (h *Handler) GetChatUsers(ctx context.Context, chatID uint32) ([]*user.User, error) {
	users, err := h.Service.GetChatUsers(ctx, chatID)
	if err != nil {
		return nil, err
	}
	return users, nil
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

// Add WebSocketHandler method to Handler
func (h *Handler) WebSocketHandler(w http.ResponseWriter, req *http.Request) {
	h.Service.WebSocketHandler(w, req)
}
