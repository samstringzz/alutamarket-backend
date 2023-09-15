package store

import (
	"github.com/Chrisentech/aluta-market-api/internals/product"
	// "github.com/Chrisentech/aluta-market-api/internals/user"
	"gorm.io/gorm"
)

type Transactions struct {
	Name string `json:"name" db:"name"`
	// Details interface{}
	StoreID uint32
}

type Product *product.Product

type Follower struct {
	FollowerID    uint32 `json:"follower_id" db:"follower_id"`
	FollowerName  string `json:"follower" db:"follower_name"`
	StoreID       uint32 `json:"store" db:"store_id"`
	FollowerImage string `json:"follower_image" db:"follower_image"`
}

type Store struct {
	gorm.Model
	ID                 uint32     `json:"id" db:"id"`
	Name               string     `json:"name" db:"name"`
	UserID             uint32     `json:"user_id" db:"user_id"`
	Products           []*Product `gorm:"foreignKey:StoreID"`
	Link               string     `json:"link" db:"link"`
	Description        string     `json:"description" db:"description"`
	HasPhysicalAddress bool       `json:"hasphysical_address" db:"has_physical_address"`
	Address            string     `json:"address" db:"address"`
	Transactions       Transactions
	Followers          []*Follower `gorm:"foreignKey:StoreID"`

	Wallet uint32 `json:"wallet" db:"wallet"`
	// Add any other fields related to a store here
}
