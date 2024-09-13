package messages

import (
	"context"
	"net/http"
	"time"

	"github.com/Chrisentech/aluta-market-api/internals/user"
	"gorm.io/gorm"
)

type Role string

const (
	Sender   Role = "sender"
	Receiver Role = "receiver"
)

type Status string

const (
	Active   Status = "active"
	Inactive Status = "inactive"
)

type MediaType string

const (
	Image MediaType = "image"
	Video MediaType = "video"
	Audio MediaType = "audio"
)

// type User struct {
// 	ID     uint32 `json:"id"`
// 	Avatar string `json:"avatar"`
// 	Role   Role   `json:"role"`
// 	Name   string `json:"name"`
// 	Status Status `json:"status" gorm:"column:status"`
// }

type Chat struct {
	gorm.Model
	ID              uint32       `json:"id" gorm:"column:id"`
	LatestMessage   Message      `gorm:"foreignKey:LatestMessageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // One-to-One or Many-to-One
	LatestMessageID *uint32      `json:"latest_message_id,omitempty"`                                                             // Foreign key for LatestMessageCreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
	CreatedAt       time.Time    `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time    `json:"updated_at" gorm:"column:updated_at"`
	UnreadCount     int          `json:"unread_count" gorm:"column:unread_count"`
	Users           []*user.User `gorm:"many2many:chat_users;" json:"chat_users"`
	Messages        []*Message   `gorm:"foreignKey:ChatID" json:"messages"` // All messages in this chat (one-to-many relationship)
}

type Message struct {
	gorm.Model
	ID        uint32       `json:"id" gorm:"column:id"`
	ChatID    uint32       `json:"chat_id" gorm:"column:chat_id"` // Foreign key to Chat
	Content   string       `json:"content" gorm:"column:content"`
	Sender    uint32       `json:"sender"` // Foreign key for LatestMessageCreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
	Media     *MediaType   `json:"media" gorm:"column:media"`
	Users     []*user.User `gorm:"serializer:json"` // Serialized users involved in this message
	CreatedAt time.Time    `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"column:updated_at"`
	IsRead    bool         `json:"is_read" gorm:"column:is_read"`
}

type Repository interface {
	SendMessage(ctx context.Context, msg *Message) error
	GetChatLists(ctx context.Context, userID uint32) ([]*Chat, error)
	FindOrCreateChat(ctx context.Context, user []*user.User) (*Chat, error)
	GetMessages(ctx context.Context, chatID string) ([]*Message, error)
	WebSocketHandler(w http.ResponseWriter, req *http.Request)
}

type Service interface {
	SendMessage(ctx context.Context, msg *Message) error
	GetChatLists(ctx context.Context, userID uint32) ([]*Chat, error)
	FindOrCreateChat(ctx context.Context, user []*user.User) (*Chat, error)
	GetMessages(ctx context.Context, chatID string) ([]*Message, error)
}
