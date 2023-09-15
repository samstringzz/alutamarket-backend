package product

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Chrisentech/aluta-market-api/errors"
	"github.com/Chrisentech/aluta-market-api/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository() Repository {
	dbURI := os.Getenv("DB_URI")

	// Initialize the database connection
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return &repository{
		db: db,
	}
}

func (r *repository) CreateCategory(ctx context.Context, req *Category) (*Category, error) {

	var count int64
	r.db.Model(&Category{}).Where("name =?", req.Name).Count(&count)
	if count > 0 {
		return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "Category already exist")
	}
	newCategory := &Category{
		Name: req.Name,
		Slug: utils.GenerateSlug(req.Name),
	}
	if err := r.db.Create(newCategory).Error; err != nil {
		return nil, err
	}
	return newCategory, nil
}

func (r *repository) CreateSubCategory(ctx context.Context, req SubCategory) (*Category, error) {

	category, err := r.GetCategory(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	req.Slug = utils.GenerateSlug(req.Name)
	category.SubCategories = append(category.SubCategories, req)
	if err := r.db.Save(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *repository) GetProduct(ctx context.Context, id uint32) (*Product, error) {
	p := Product{}
	err := r.db.Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) CreateProduct(ctx context.Context, req *Product) (*Product, error) {

	newProduct := &Product{
		Name:          req.Name,
		Slug:          utils.GenerateSlug(req.Name),
		Description:   req.Description,
		Image:         req.Image,
		Price:         req.Price,
		Status:        req.Status,
		Quantity:      req.Quantity,
		Campus:        req.Campus,
		Variant:       req.Variant,
		Condition:     req.Condition,
		StoreID:       req.StoreID,
		CategoryID:    req.CategoryID,
		SubCategoryID: req.SubCategoryID,
	}
	if err := r.db.Create(newProduct).Error; err != nil {

		log.Printf("Error creating product: %v", err)
		return nil, err
	}
	return newProduct, nil
}

// func (r *repository) UpdateProduct(ctx context.Context, req *Product)(*Product,error){}
// func (r *repository) DeleteProduct(ctx context.Context, id int)(error){}
// func (r *repository) GetProduct(ctx context.Context, filter string,filterOption string)(*Product,error){}
func (r *repository) GetProducts(ctx context.Context)([]*Product,error){
		var products []*Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
func (r *repository) GetCategories(ctx context.Context) ([]*Category, error) {
	var categories []*Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
func (r *repository) GetCategory(ctx context.Context, id uint32) (*Category, error) {
	p := Category{}
	err := r.db.Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// func (r *repository) RecommendedProducts(ctx context.Context, userId uint32, productId uint32)([]*RecommendedProduct,error){
// cat:= Category{}
// prd:= Product{}
// product:= &RecommendedProduct{}
// store:=Store{}
// if err :=  r.db.Where("id = ?", productId).First(&prd).Error; err != nil {
// 		return nil, err
// 	}
// allProducts,_ := r.GetProducts(ctx)

// for _, p:= range allProducts{
// []product = ([]product,p)
// }
// }