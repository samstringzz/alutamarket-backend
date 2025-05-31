package models

import "time"

type Downloads struct {
	ID        uint32    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id"`
	ProductID uint32    `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
