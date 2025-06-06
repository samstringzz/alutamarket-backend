package subscriber

import (
	"errors"
	"fmt"

	"github.com/samstringzz/alutamarket-backend/database"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository() Repository {
	return &repository{
		db: database.GetDB(),
	}
}

func (r *repository) CreateSubscriber(email string) (*Subscriber, error) {
	// First check if email already exists
	var existingSubscriber Subscriber
	if err := r.db.Where("email = ?", email).First(&existingSubscriber).Error; err == nil {
		// Email already exists
		return nil, fmt.Errorf("email %s is already subscribed to our newsletter", email)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Some other database error occurred
		return nil, err
	}

	subscriber := &Subscriber{
		Email:  email,
		Active: true,
	}

	if err := r.db.Create(subscriber).Error; err != nil {
		return nil, fmt.Errorf("failed to create subscriber: %v", err)
	}

	return subscriber, nil
}

func (r *repository) GetSubscribers() ([]*Subscriber, error) {
	var subscribers []*Subscriber
	if err := r.db.Where("active = ?", true).Find(&subscribers).Error; err != nil {
		return nil, err
	}
	return subscribers, nil
}

func (r *repository) UnsubscribeEmail(email string) error {
	return r.db.Model(&Subscriber{}).Where("email = ?", email).Update("active", false).Error
}
