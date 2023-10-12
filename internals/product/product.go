package product

import (
	"context"

	"gorm.io/gorm"
	// "github.com/Chrisentech/aluta-market-api/internals/store"
)

type Category struct {
	gorm.Model
	ID            int           `json:"id" db:"id"`
	Name          string        `json:"name" db:"name"`
	Slug          string        `json:"slug" db:"slug"`
	SubCategories []SubCategory `gorm:"serializer:json"`
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
type NewProduct struct {
	gorm.Model
	ID            uint32   `json:"id" db:"id"`
	Name          string   `json:"name" db:"name"`
	Description   string   `json:"description" db:"description"`
	Image         []string `gorm:"serializer:json" json:"image" db:"image"`
	Thumbnail     string   `json:"thumbnail" db:"thumbnail"`
	Price         float64  `json:"price" db:"price"`
	Discount      float64  `json:"discount" db:"discount"`
	Status        bool     `json:"status" db:"status"`
	Quantity      int      `json:"quantity" db:"quantity"`
	Slug          string   `json:"slug" db:"slug"`
	Variant       string   `json:"variant,omitempty" db:"variant"`
	Store         string   `json:"store" db:"store_id"`
	CategoryID    uint8    `json:"category" db:"category_id"`
	SubCategoryID uint8    `json:"subcategory" db:"sub_category_id"`
}
type Product struct {
	gorm.Model
	ID          uint32    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Image       []string  `gorm:"serializer:json" json:"image" db:"image"`
	Thumbnail   string    `json:"thumbnail" db:"thumbnail"`
	Price       float64   `json:"price" db:"price"`
	Discount    float64   `json:"discount" db:"discount"`
	Status      bool      `json:"status" db:"status"`
	Quantity    int       `json:"quantity" db:"quantity"`
	Slug        string    `json:"slug" db:"slug"`
	Variant     string    `json:"variant,omitempty" db:"variant"`
	Store       string    `json:"store" db:"store"`
	Category    string    `json:"category" db:"category"`
	Views       []uint32  `gorm:"serializer:json" jsinput.ProductIDon:"views" db:"views"`
	Subcategory string    `json:"subcategory" db:"subcategory"`
	Reviews     []*Review `gorm:"serializer:json"`
}

type WishListedProduct struct {
	gorm.Model
	UserID  uint32   `json:"user_id" db:"user_id"`
	Product *Product `gorm:"embedded"`
}

type RecommendedProduct struct {
	gorm.Model
	UserID  uint32   `json:"user_id" db:"user_id"`
	Product *Product `gorm:"embedded,type:products"`
}

type Repository interface {
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	CreateSubCategory(ctx context.Context, subcategory SubCategory) (*Category, error)
	GetCategories(ctx context.Context) ([]*Category, error)
	GetCategory(ctx context.Context, id uint32) (*Category, error)
	CreateProduct(ctx context.Context, product *NewProduct) (*Product, error)
	GetProduct(ctx context.Context, id uint32) (*Product, error)
	GetProducts(ctx context.Context, store string, limit int, offset int) ([]*Product, error)
	AddWishListedProduct(ctx context.Context, userId, productId uint32) (*WishListedProduct, error)
	GetWishListedProducts(ctx context.Context, userId uint32) ([]*WishListedProduct, error)
	GetRecommendedProducts(ctx context.Context, query string)([]*Product,error)
	SearchProducts(ctx context.Context, query string)([]*Product,error)
	RemoveWishListedProduct(ctx context.Context, userId uint32) error
	// GetProductByFilter(ctx context.Context, filter string,filterOption string )(*Product,error)    //by slug,by store,by id,(by category||subcategory)
	UpdateProduct(ctx context.Context, req *Product) (*Product, error)
	DeleteProduct(ctx context.Context, id uint32)(error)

}

type Service interface {
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	CreateSubCategory(ctx context.Context, subcategory SubCategory) (*Category, error)
	GetCategories(ctx context.Context) ([]*Category, error)
	GetCategory(ctx context.Context, id uint32) (*Category, error)
	CreateProduct(ctx context.Context, product *NewProduct) (*Product, error)
	GetProducts(ctx context.Context, store string, limit int, offset int) ([]*Product, error)
	GetProduct(ctx context.Context, id uint32) (*Product, error)
	GetRecommendedProducts(ctx context.Context, query string)([]*Product,error)
	// GetProductByFilter(ctx context.Context, filter string,filterOption string)(*Product,error)    //by slug,by store,by id,(by category||subcategory)
	SearchProducts(ctx context.Context, query string)([]*Product,error)
	UpdateProduct(ctx context.Context, req *Product) (*Product, error)
	AddWishListedProduct(ctx context.Context, userId, productId uint32) (*WishListedProduct, error)
	GetWishListedProducts(ctx context.Context, userId uint32) ([]*WishListedProduct, error)
	RemoveWishListedProduct(ctx context.Context, userId uint32) error
	DeleteProduct(ctx context.Context, id uint32 )(error)

	// Left are
	// Add Review to Product/Store?? yet to be decided

}
