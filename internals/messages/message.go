package messages

import (
	"context"
	"time"

	"github.com/lib/pq"
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

type User struct {
	ID     uint32 `json:"id"`
	Avatar string `json:"avatar"`
	Role   Role   `json:"role"`
	Name   string `json:"name"`
	Status Status `json:"status" gorm:"column:status"`
}

type Chat struct {
	gorm.Model
	ID            uint32        `json:"id" gorm:"column:id"`
	LatestMessage Message       `json:"latest_message" gorm:"embedded"`
	CreatedAt     time.Time     `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"column:updated_at"`
	UnreadCount   int           `json:"unread_count" gorm:"column:unread_count"`
	Users         pq.Int64Array `json:"users" gorm:"type:bigint[]"` // Correct type
}

type Message struct {
	gorm.Model
	ID        uint32    `json:"id" gorm:"column:id"`
	ChatID    uint32    `json:"chat_id" gorm:"column:chat_id"`
	Content   string    `json:"content" gorm:"column:content"`
	Media     MediaType `json:"media" gorm:"column:media"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	User      User      `json:"user" gorm:"embedded"`
	IsRead    bool      `json:"is_read" gorm:"column:is_read"`
}

type Repository interface {
	SendMessage(ctx context.Context, msg *Message) error
	GetChatLists(ctx context.Context, userID int64) ([]*Chat, error)
	CreateChat(ctx context.Context, chatID uint32, usersID pq.Int64Array) error
	GetMessages(ctx context.Context, chatID string) ([]*Message, error)
}

type Service interface {
	SendMessage(ctx context.Context, msg *Message) error
	GetChatLists(ctx context.Context, userID int64) ([]*Chat, error)
	CreateChat(ctx context.Context, chatID uint32, usersID pq.Int64Array) error
	GetMessages(ctx context.Context, chatID string) ([]*Message, error)
}
