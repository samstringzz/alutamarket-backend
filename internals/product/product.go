//Ok HandledProducts would hande wishlist,savedforLater,recentlyAdded

package product

import (
	"context"
	"time"
	"gorm.io/gorm"
)

type AdsGen struct {
	gorm.Model
	ID       uint32    `json:"id" db:"id"`
	Units    uint8     `json:"units" db:"units"`
	Validity time.Time `json:"validity" db:"validity"`
}
type Category struct {
	gorm.Model
	ID            int           `json:"id" db:"id"`
	Name          string        `json:"name" db:"name"`
	Slug          string        `json:"slug" db:"slug"`
	SubCategories []SubCategory `gorm:"serializer:json"`
}
type SubCategory struct {
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
	Rating    float64  `json:"rating" db:"rating"`
	ProductID uint32 `json:"product" db:"product_id"`
}
type VariantValue struct {
	Value  string   `json:"value" db:"value"`
	Price  float64  `json:"price,omitempty" db:"price"`
	Images []string `gorm:"serializer:json"  json:"images,omitempty" db:"images"`
}
type VariantType struct {
	Name  string          `json:"variant_name" db:"varaint_name"`
	Value []*VariantValue `gorm:"serializer:json" json:"variant_value" db:"variant_vaue"`
}
type NewProduct struct {
	gorm.Model
	ID            uint32         `json:"id" db:"id"`
	Name          string         `json:"name" db:"name"`
	Description   string         `json:"description" db:"description"`
	Images        []string       `gorm:"serializer:json" json:"image" db:"image"`
	Thumbnail     string         `json:"thumbnail" db:"thumbnail"`
	Price         float64        `json:"price" db:"price"`
	Discount      float64        `json:"discount" db:"discount"`
	Status        bool           `json:"status" db:"status"`
	Quantity      int            `json:"quantity" db:"quantity"`
	Slug          string         `json:"slug" db:"slug"`
	Variant       []*VariantType `gorm:"serializer:json" json:"variant,omitempty" db:"variant"`
	Store         string         `json:"store" db:"store"`
	CategoryID    uint8          `json:"category" db:"category_id"`
	SubCategoryID uint8          `json:"subcategory" db:"sub_category_id"`
}
type Product struct {
	gorm.Model
	ID          uint32         `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Description string         `json:"description" db:"description"`
	Images      []string       `gorm:"serializer:json" json:"image" db:"image"`
	Thumbnail   string         `json:"thumbnail" db:"thumbnail"`
	Price       float64        `json:"price" db:"price"`
	Discount    float64        `json:"discount" db:"discount"`
	Status      bool           `json:"status" db:"status"`
	Quantity    int            `json:"quantity" db:"quantity"`
	Slug        string         `json:"slug" db:"slug"`
	Variant     []*VariantType `gorm:"serializer:json" json:"variant,omitempty" db:"variant"`
	Store       string         `json:"store" db:"store"`
	Category    string         `json:"category" db:"category"`
	Views       []uint32       `gorm:"serializer:json" jsinput.ProductIDon:"views" db:"views"`
	Subcategory string         `json:"subcategory" db:"subcategory"`
	Reviews     []*Review      `gorm:"serializer:json"`
	Ads         *AdsGen        `gorm:"serializer:json" json:"ads,omitempty" db:"ads"`
}

type WishListedProduct struct {
	gorm.Model
	UserID  uint32   `json:"user_id" db:"user_id"`
	Product *Product `gorm:"embedded"`
}

type RecentlyViewedProduct struct {
	gorm.Model
	UserID  uint32   `json:"user_id" db:"user_id"`
	Product *Product `gorm:"embedded"`
}

type HandledProduct struct {
	gorm.Model
	UserID  uint32   `json:"user_id" db:"user_id"`
	Product *Product `gorm:"embedded"`
    Type    string   `db:"type"`
}

type Repository interface {
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	CreateSubCategory(ctx context.Context, subcategory SubCategory) (*Category, error)
	GetCategories(ctx context.Context) ([]*Category, error)
	GetCategory(ctx context.Context, id uint32) (*Category, error)
	CreateProduct(ctx context.Context, product *NewProduct) (*Product, error)
	GetProduct(ctx context.Context, productId, userId uint32) (*Product, error)
	GetProducts(ctx context.Context, store string, limit int, offset int) ([]*Product, error)
	AddHandledProduct(ctx context.Context, userId, productId uint32,eventType string) (*HandledProduct, error)
	AddSavedForLater(ctx context.Context, userId, productId uint32) (*HandledProduct, error)
	GetHandledProducts(ctx context.Context, userId uint32,eventType string) ([]*HandledProduct, error)
	GetRecommendedProducts(ctx context.Context, query string) ([]*Product, error)
	SearchProducts(ctx context.Context, query string) ([]*Product, error)
	RemoveHandledProduct(ctx context.Context, userId uint32,eventType string) error
	// GetProductByFilter(ctx context.Context, filter string,filterOption string )(*Product,error)    //by slug,by store,by id,(by category||subcategory)
	UpdateProduct(ctx context.Context, req *Product) (*Product, error)
	DeleteProduct(ctx context.Context, id uint32) error
	// GetRecentlyViewedProducts(ctx context.Context, userId uint32) ([]*Product, error)
	// AddRecentlyViewedProducts(ctx context.Context, userId,productId uint32) error
	AddReview(ctx context.Context, input *Review) (*Review,error)
	GetReviews(ctx context.Context, productId uint32) ([]*Review,error)
}

type Service interface {
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	CreateSubCategory(ctx context.Context, subcategory SubCategory) (*Category, error)
	GetCategories(ctx context.Context) ([]*Category, error)
	GetCategory(ctx context.Context, id uint32) (*Category, error)
	CreateProduct(ctx context.Context, product *NewProduct) (*Product, error)
	AddHandledProduct(ctx context.Context, userId, productId uint32,eventType string) (*HandledProduct, error)
	GetProducts(ctx context.Context, store string, limit int, offset int) ([]*Product, error)
	GetProduct(ctx context.Context, productId, userId uint32) (*Product, error)
	GetRecommendedProducts(ctx context.Context, query string) ([]*Product, error)
	// GetProductByFilter(ctx context.Context, filter string,filterOption string)(*Product,error)    //by slug,by store,by id,(by category||subcategory)
	GetHandledProducts(ctx context.Context, userId uint32,eventType string) ([]*HandledProduct, error)	
	SearchProducts(ctx context.Context, query string) ([]*Product, error)
	UpdateProduct(ctx context.Context, req *Product) (*Product, error)
	// AddWishListedProduct(ctx context.Context, userId, productId uint32) (*HandledProduct, error)
	// GetWishListedProducts(ctx context.Context, userId uint32) ([]*HandledProduct, error)
	RemoveHandledProduct(ctx context.Context, userId uint32,eventType string) error
	// RemoveWishListedProduct(ctx context.Context, userId uint32) error
	// AddRecentlyViewedProducts(ctx context.Context, userId,productId uint32) error
	// GetRecentlyViewedProducts(ctx context.Context, userId uint32) ([]*Product, error)
	AddReview(ctx context.Context, input *Review) (*Review,error)
	GetReviews(ctx context.Context, productId uint32) ([]*Review,error)
	DeleteProduct(ctx context.Context, id uint32) error

	// Left are
	// Add Review to Product??
	//

}
