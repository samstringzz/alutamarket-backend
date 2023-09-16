package product

import (
	"context"

	"gorm.io/gorm"
	// "github.com/Chrisentech/aluta-market-api/internals/store"
)

// type Store *store.Store

type Category struct {
	gorm.Model
	ID            int           `json:"id" db:"id"`
	Name          string        `json:"name" db:"name"`
	Slug          string        `json:"slug" db:"slug"`
	SubCategories []SubCategory `gorm:"foreignKey:CategoryID"`
}
type SubCategory struct {
	gorm.Model
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	CategoryID uint32 `json:"category"`
}
type Order struct{}

type Review struct {
	gorm.Model
	Username  string `json:"username" db:"username"`
	Image     string `json:"image" db:"image"`
	Message   string `json:"message" db:"message"`
	Rating    uint8  `json:"rating" db:"rating"`
	ProductID uint32 `json:"product" db:"product_id"`
}

type Product struct {
	gorm.Model
	ID            uint32    `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	Image         []string  `json:"image" db:"image"`
	Thumbnail     string    `json:"thumbnail" db:"thumbnail"`
	Price         float64   `json:"price" db:"price"`
	Discount      float64   `json:"discount" db:"discount"`
	Status        bool      `json:"status" db:"status"`
	Quantity      int       `json:"quantity" db:"quantity"`
	Slug          string    `json:"slug" db:"slug"`
	Variant       string    `json:"variant,omitempty" db:"variant"`
	StoreID       uint32    `json:"store" db:"store_id"`
	CategoryID    uint8     `json:"category" db:"category_id"`
	Views         []*uint32 `gorm:"foreignKey:UserID"`
	SubCategoryID uint8     `json:"subcategory" db:"sub_category_id"`
	Reviews       []*Review `gorm:"foreignkey:ProductID"`
}

type WishListedProduct struct {
	gorm.Model
	UserID  uint32 `json:"user_id" db:"user_id"`
	Product *Product `gorm:"embedded"`
}

type RecommendedProduct struct {
	gorm.Model
	UserID  uint32 `json:"user_id" db:"user_id"`
	Product *Product `gorm:"embedded"`
}

type Repository interface {
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	CreateSubCategory(ctx context.Context, subcategory SubCategory) (*Category, error)
	GetCategories(ctx context.Context) ([]*Category, error)
	GetCategory(ctx context.Context, id uint32) (*Category, error)
	CreateProduct(ctx context.Context, product *Product) (*Product, error)
	GetProduct(ctx context.Context, id uint32) (*Product, error)
	GetProducts(ctx context.Context, store string) ([]*Product, error)
	AddWishListedProduct(ctx context.Context, userId,productId uint32) (*WishListedProduct,error)
	GetWishListedProducts(ctx context.Context, userId uint32) ([]*WishListedProduct,error)
	RemoveWishListedProduct(ctx context.Context, userId uint32) error
	// GetProductByFilter(ctx context.Context, filter string,filterOption string )(*Product,error)    //by slug,by store,by id,(by category||subcategory)
	UpdateProduct(ctx context.Context, req *Product) (*Product, error)
	// DeleteProduct(ctx context.Context, id int)(error)
	//Recommended Products
	// Wishlisted Product

}

type Service interface {
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	CreateSubCategory(ctx context.Context, subcategory SubCategory) (*Category, error)
	GetCategories(ctx context.Context) ([]*Category, error)
	GetCategory(ctx context.Context, id uint32) (*Category, error)
	CreateProduct(ctx context.Context, product *Product) (*Product, error)
	GetProducts(ctx context.Context, store string) ([]*Product, error)
	GetProduct(ctx context.Context, id uint32) (*Product, error)
	// GetProductByFilter(ctx context.Context, filter string,filterOption string)(*Product,error)    //by slug,by store,by id,(by category||subcategory)
	UpdateProduct(ctx context.Context, req *Product) (*Product, error)
	AddWishListedProduct(ctx context.Context, userId,productId uint32) (*WishListedProduct,error)
	GetWishListedProducts(ctx context.Context, userId uint32) ([]*WishListedProduct,error)
	RemoveWishListedProduct(ctx context.Context, userId uint32) error
	// DeleteProduct(ctx context.Context, id int )(error)

	// Left are
	//Recommended Product
	//Add to wishlist and get wishlisted product
	// Add Review to Product/Store?? yet to be decided

}
