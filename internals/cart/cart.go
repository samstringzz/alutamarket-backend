package cart

import (
	"context"
	"net/http"

	"github.com/Chrisentech/aluta-market-api/internals/product"
	"github.com/Chrisentech/aluta-market-api/internals/store"
	"github.com/Chrisentech/aluta-market-api/internals/user"
	"gorm.io/gorm"
)

type Product *product.Product
type User *user.User
type Order *store.Order
type Cart struct {
	gorm.Model
	ID       uint32       `gorm:"primaryKey;uniqueIndex;not null"`
	Items    []*CartItems `gorm:"serializer:json" json:"items" db:"items"`
	Total    float64      `json:"total" db:"total"`
	Active   bool         `json:"active" db:"active"`
	UserID   uint32       `json:"user" db:"user_id"`
	StoresID []*string    `gorm:"serializer:json" json:"stores" db:"stores"`
}

type CartItems struct {
	Product  Product `gorm:"embedded"`
	CartID   uint32  `json:"cart" db:"cart_id"`
	Quantity int     `json:"quantity" db:"quantity"`
	Fee      float64 `json:"fee" db:"fee"`
}

type Repository interface {
	ModifyCart(ctx context.Context, req *CartItems, user uint32) (*Cart, error)
	RemoveAllCart(ctx context.Context, id uint32) error
	GetCart(ctx context.Context, user uint32) (*Cart, error)
	MakePayment(ctx context.Context, w http.ResponseWriter, r *http.Request)
	InitiatePayment(ctx context.Context, req Order) (string, error)
}

type Service interface {
	ModifyCart(ctx context.Context, req *CartItems, user uint32) (*Cart, error)
	RemoveAllCart(ctx context.Context, id uint32) error
	GetCart(ctx context.Context, user uint32) (*Cart, error)
	MakePayment(ctx context.Context, w http.ResponseWriter, r *http.Request)
	InitiatePayment(ctx context.Context, req Order) (string, error)
	GetProduct(ctx context.Context, productId uint32) (*product.Product, error)
}
