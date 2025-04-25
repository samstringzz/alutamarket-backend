package subscriber

import (
    "time"
    "gorm.io/gorm"
)

type Subscriber struct {
    gorm.Model
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    Active    bool      `json:"active" gorm:"default:true"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
    CreateSubscriber(email string) (*Subscriber, error)
    GetSubscribers() ([]*Subscriber, error)
}

type Service interface {
    CreateSubscriber(email string) (*Subscriber, error)
    GetSubscribers() ([]*Subscriber, error)
}