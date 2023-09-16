package product

import (
	"context"
	"fmt"
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
		Thumbnail:     req.Thumbnail,
		Price:         req.Price,
		Discount:      req.Discount,
		Status:        req.Status,
		Quantity:      req.Quantity,
		Variant:       req.Variant,
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

func (r *repository) UpdateProduct(ctx context.Context, req *Product) (*Product, error) {

	// First, check if the product exists by its ID or another unique identifier
	existingProduct, err := r.GetProduct(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// Update only the fields that are present in the req
	if req.Name != "" {
		existingProduct.Name = req.Name
		existingProduct.Slug = utils.GenerateSlug(req.Name)
	}
	if req.Description != "" {
		existingProduct.Description = req.Description
	}
	if req.Quantity != 0 {
		existingProduct.Quantity = req.Quantity
	}
	if len(req.Image) != 0 {
		existingProduct.Image = append(existingProduct.Image, req.Image...)
	}
	if req.Discount != 0 {
		existingProduct.Discount = req.Discount
	}
	if req.Status {
		existingProduct.Status = req.Status
	}
	if req.Thumbnail != "" {
		existingProduct.Thumbnail = req.Thumbnail
	}
	if req.Price != 0 {
		existingProduct.Price = req.Price
	}

	if req.Variant != "" {
		existingProduct.Variant = req.Variant
	}

	// Update the product in the repository
	err = r.db.Save(existingProduct).Error
	if err != nil {
		return nil, err
	}

	return existingProduct, nil
}

func (r *repository) DeleteProduct(ctx context.Context, id uint32) error {
	existingProduct, err := r.GetProduct(ctx, id)
	if err != nil {
		return err
	}
	err = r.db.Delete(existingProduct).Error
	return err

}
func (r *repository) AddWishListedProduct(ctx context.Context, userId, productId uint32) (*WishListedProduct, error) {
	wishlist := &WishListedProduct{}
	foundProduct, err := r.GetProduct(ctx, productId)
	if err != nil {
		return nil, err
	}
	var count int64
	r.db.Model(wishlist).Where("user_id =?", userId).Count(&count)
	if count > 0 {
		fmt.Printf("The Total no of User wishlist is%v\n", count)
		return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "Product already in wishlist")
	}
	wishlist.Product = foundProduct
	wishlist.UserID = userId
	err = r.db.Create(wishlist).Error
	if err != nil {
		return nil, err
	}
	return wishlist, nil
}

func (r *repository) GetWishListedProducts(ctx context.Context, userId uint32) ([]*WishListedProduct, error) {
	var wishlist []*WishListedProduct
	if err := r.db.Find(&wishlist).Where("user_id",userId).Error; err != nil {
		return nil, err
	}
	return wishlist,nil
}

func (r *repository) RemoveWishListedProduct(ctx context.Context, userId uint32) error {
	existingWishlist :=&WishListedProduct{}
	err := r.db.Delete(existingWishlist).Error
	return err
}

// recommended and wishlisted Product
// func (r *repository) GetProduct(ctx context.Context, filter string,filterOption string)(*Product,error){}
func (r *repository) GetProducts(ctx context.Context, store string) ([]*Product, error) {
	var products []*Product
	if store != "" {
		if err := r.db.Where("store_id = ? ", store).Find(&products).Error; err != nil {
			return nil, err
		}
		return products, nil
	} else {
		if err := r.db.Where("status = ? ", true).Find(&products).Error; err != nil {
			return nil, err
		}
		return products, nil
	}

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
