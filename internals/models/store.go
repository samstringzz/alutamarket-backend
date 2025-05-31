package models

import (
	"time"
)

type Follower struct {
	FollowerID    uint32 `json:"follower_id"`
	FollowerName  string `json:"follower_name"`
	FollowerImage string `json:"follower_image"`
}

type Store struct {
	ID                 uint32      `json:"id" gorm:"primaryKey"`
	Name               string      `json:"name"`
	Email              string      `json:"email"`
	Link               string      `json:"link"`
	UserID             uint32      `json:"user_id"`
	Description        string      `json:"description"`
	HasPhysicalAddress bool        `json:"has_physical_address"`
	Address            string      `json:"address"`
	Wallet             float64     `json:"wallet"`
	Status             bool        `json:"status"`
	Background         string      `json:"background"`
	Phone              string      `json:"phone"`
	Thumbnail          string      `json:"thumbnail"`
	Followers          []*Follower `json:"followers" gorm:"type:jsonb"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
}
