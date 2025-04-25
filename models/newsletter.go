package models

import "time"

type NewsletterSubscriber struct {
    ID           uint      `gorm:"primaryKey"`
    Email        string    `gorm:"unique;not null"`
    SubscribedAt time.Time `gorm:"not null"`
}