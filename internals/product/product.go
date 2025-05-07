//Ok HandledProducts would hande wishlist,savedforLater,recentlyAdded

package product

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Chrisentech/aluta-market-api/graph/model"
	"github.com/lib/pq"
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
	Type          string        `json:"type" db:"type"`
	SubCategories []SubCategory `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE;" json:"subcategories"`
}

type SubCategory struct {
	gorm.Model
	Name       string     `json:"name"`
	Slug       string     `json:"slug"`
	CategoryID uint32     `json:"category_id"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

type Order struct{}

type Review struct {
	gorm.Model
	ID        string  `json:"id" db:"id"`
	Username  string  `json:"username" db:"username"`
	Image     string  `json:"image" db:"image"`
	Message   string  `json:"message" db:"message"`
	Seller    string  `json:"seller" db:"seller"`
	Rating    float64 `json:"rating" db:"rating"`
	ProductID uint32  `json:"product" db:"product_id"`
}
type VariantValue struct {
	Value  string   `json:"value" db:"value"`
	Price  float64  `json:"price,omitempty" db:"price"`
	Images []string `gorm:"serializer:json"  json:"images,omitempty" db:"images"`
}
type VariantType struct {
	Name  string          `json:"variant_name" db:"varaint_name"`
	Price *float64        `json:"variant_price" db:"varaint_price"`
	Value []*VariantValue `gorm:"serializer:json" json:"variant_value" db:"variant_value"`
}

type NewProduct struct {
	gorm.Model
	ID              string         `json:"id" db:"id"`
	Name            string         `json:"name" db:"name"`
	Description     string         `json:"description" db:"description"`
	Images          []string       `gorm:"serializer:json" json:"image" db:"image"`
	Thumbnail       string         `json:"thumbnail" db:"thumbnail"`
	Price           float64        `json:"price" db:"price"`
	Discount        float64        `json:"discount" db:"discount"`
	Status          *bool          `json:"status" db:"status"`
	Quantity        int            `json:"quantity" db:"quantity"`
	File            string         `json:"file" db:"file"`
	Slug            string         `json:"slug" db:"slug"`
	Variant         []*VariantType `gorm:"serializer:json" json:"variant,omitempty" db:"variant"`
	Store           string         `json:"store" db:"store"`
	CategoryID      uint8          `json:"category" db:"category_id"`
	SubCategoryName string         `json:"subcategory" db:"sub_category_name"`
	AlwaysAvailbale bool           `json:"always_available" db:"always_available"`
}

type UpdateStore struct {
	gorm.Model
	ID              uint32         `json:"id" db:"id"`
	Name            string         `json:"name" db:"name"`
	Description     string         `json:"description" db:"description"`
	Images          []string       `gorm:"serializer:json" json:"image" db:"image"`
	Thumbnail       string         `json:"thumbnail" db:"thumbnail"`
	Price           float64        `json:"price" db:"price"`
	Discount        float64        `json:"discount" db:"discount"`
	Status          bool           `json:"status" db:"status"`
	Quantity        int            `json:"quantity" db:"quantity"`
	File            string         `json:"file" db:"file"`
	Slug            string         `json:"slug" db:"slug"`
	Variant         []*VariantType `gorm:"serializer:json" json:"variant,omitempty" db:"variant"`
	Store           string         `json:"store" db:"store"`
	CategoryID      uint8          `json:"category" db:"category_id"`
	SubCategoryID   uint8          `json:"subcategory" db:"sub_category_id"`
	AlwaysAvailbale bool           `json:"always_available" db:"always_available"`
}

// Add this type and methods at the top of the file, after imports
type Uint32Array []uint32

func (a *Uint32Array) Scan(value interface{}) error {
	if value == nil {
		*a = Uint32Array{}
		return nil
	}

	switch v := value.(type) {
	case string:
		// Remove { and } from PostgreSQL array string
		trimmed := strings.Trim(v, "{}")
		if trimmed == "" {
			*a = Uint32Array{}
			return nil
		}

		// Split the string into individual numbers
		nums := strings.Split(trimmed, ",")
		result := make([]uint32, 0, len(nums))

		for _, num := range nums {
			n, err := strconv.ParseUint(strings.TrimSpace(num), 10, 32)
			if err != nil {
				return err
			}
			result = append(result, uint32(n))
		}
		*a = result
		return nil
	default:
		return fmt.Errorf("unsupported Scan type for Uint32Array: %T", value)
	}
}

func (a Uint32Array) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}

	var b strings.Builder
	b.WriteString("{")
	for i, v := range a {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(strconv.FormatUint(uint64(v), 10))
	}
	b.WriteString("}")
	return b.String(), nil
}

// Update the Product struct
type Product struct {
	ID              uint32         `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name"`
	Slug            string         `json:"slug"`
	Description     string         `json:"description"`
	Images          pq.StringArray `json:"images" gorm:"type:text[]"`
	Thumbnail       string         `json:"thumbnail"`
	Price           float64        `json:"price"`
	Discount        float64        `json:"discount"`
	Status          bool           `json:"status"`
	AlwaysAvailbale bool           `json:"always_available"`
	Quantity        int            `json:"quantity"`
	File            string         `json:"file"`
	Store           string         `json:"store"`
	Category        string         `json:"category"`
	Subcategory     string         `json:"subcategory"`
	Type            string         `json:"type"`
	UnitsSold       int            `json:"units_sold"`
	Variant         []*VariantType `json:"variant" gorm:"type:jsonb;serializer:json"`
	Views           Uint32Array    `json:"views" gorm:"type:integer[]"`
	Reviews         []Review       `json:"reviews" gorm:"type:jsonb;serializer:json"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type ProductPaginationData struct {
	Data        []*Product
	CurrentPage int
	PerPage     int
	Total       int
	NextPage    int
	PrevPage    int
}

type HandledProduct struct {
	gorm.Model
	UserID          uint32         `json:"user_id" gorm:"column:user_id"`
	ProductID       uint32         `json:"product_id" gorm:"column:product_id"`
	Type            string         `json:"type" gorm:"column:type"`
	Name            string         `json:"name" gorm:"column:name"`
	Slug            string         `json:"slug" gorm:"column:slug"`
	Description     string         `json:"description" gorm:"column:description"`
	Images          pq.StringArray `json:"images" gorm:"type:text[]"`
	Thumbnail       string         `json:"thumbnail" gorm:"column:thumbnail"`
	Price           float64        `json:"price" gorm:"column:price"`
	Discount        float64        `json:"discount" gorm:"column:discount"`
	Status          bool           `json:"status" gorm:"column:status"`
	Quantity        int            `json:"quantity" gorm:"column:quantity"`
	File            string         `json:"file" gorm:"column:file"`
	Store           string         `json:"store" gorm:"column:store"`
	Category        string         `json:"category" gorm:"column:category"`
	Subcategory     string         `json:"subcategory" gorm:"column:subcategory"`
	UnitSold        int            `json:"unit_sold" gorm:"column:unit_sold"`
	Variant         []*VariantType `json:"variant" gorm:"type:jsonb;serializer:json"`
	Views           Uint32Array    `json:"views" gorm:"type:integer[]"`
	Reviews         []Review       `json:"reviews" gorm:"type:jsonb;serializer:json"`
	AlwaysAvailbale bool           `json:"always_availbale" gorm:"column:always_availbale"`
	Product         *Product       `json:"-" gorm:"foreignKey:ProductID"`
}

type Repository interface {
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	CreateSubCategory(ctx context.Context, subcategory SubCategory) (*Category, error)
	GetCategories(ctx context.Context) ([]*Category, error)
	GetCategory(ctx context.Context, id uint32) (*Category, error)
	CreateProduct(ctx context.Context, product *NewProduct) (*Product, error)
	GetProduct(ctx context.Context, productId, userId uint32) (*Product, error)
	GetProducts(ctx context.Context, store string, categorySlug string, limit int, offset int) ([]*Product, int, error)
	AddHandledProduct(ctx context.Context, userId, productId uint32, eventType string) (*HandledProduct, error)
	AddSavedForLater(ctx context.Context, userId, productId uint32) (*HandledProduct, error)
	GetHandledProducts(ctx context.Context, userId uint32, eventType string) ([]*HandledProduct, error)
	GetRecommendedProducts(ctx context.Context, query string) ([]*Product, error)
	SearchProducts(ctx context.Context, query string) ([]*Product, error)
	RemoveHandledProduct(ctx context.Context, userId uint32, eventType string) error
	// GetProductByFilter(ctx context.Context, filter string,filterOption string )(*Product,error)    //by slug,by store,by id,(by category||subcategory)
	UpdateProduct(ctx context.Context, req *NewProduct) (*Product, error)
	DeleteProduct(ctx context.Context, id uint32) error
	AddReview(ctx context.Context, input *Review) (*Review, error)
	GetProductReviews(ctx context.Context, productId uint32, sellerId string) ([]*Review, error)
	Products(ctx context.Context, store *string, categorySlug *string, limit *int, offset *int) (*model.ProductPaginationData, error)
}

type Service interface {
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	CreateSubCategory(ctx context.Context, subcategory SubCategory) (*Category, error)
	GetCategories(ctx context.Context) ([]*Category, error)
	GetCategory(ctx context.Context, id uint32) (*Category, error)
	CreateProduct(ctx context.Context, product *NewProduct) (*Product, error)
	AddHandledProduct(ctx context.Context, userId, productId uint32, eventType string) (*HandledProduct, error)
	GetProducts(ctx context.Context, store string, categorySlug string, limit int, offset int) ([]*Product, int, error)
	GetProduct(ctx context.Context, productId, userId uint32) (*Product, error)
	GetRecommendedProducts(ctx context.Context, query string) ([]*Product, error)
	// GetProductByFilter(ctx context.Context, filter string,filterOption string)(*Product,error)    //by slug,by store,by id,(by category||subcategory)
	GetHandledProducts(ctx context.Context, userId uint32, eventType string) ([]*HandledProduct, error)
	SearchProducts(ctx context.Context, query string) ([]*Product, error)
	UpdateProduct(ctx context.Context, req *NewProduct) (*Product, error)
	RemoveHandledProduct(ctx context.Context, userId uint32, eventType string) error
	AddReview(ctx context.Context, input *Review) (*Review, error)
	GetProductReviews(ctx context.Context, productId uint32, sellerId string) ([]*Review, error)
	DeleteProduct(ctx context.Context, id uint32) error
}
