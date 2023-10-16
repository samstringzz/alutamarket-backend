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

func (r *repository) CreateProduct(ctx context.Context, req *NewProduct) (*Product, error) {
	category, _ := r.GetCategory(ctx, uint32(req.CategoryID))
	subcategory := req.SubCategoryID
	subcategoryName := ""
	for i, item := range category.SubCategories {
		if i+1 == int(subcategory) {
			subcategoryName = item.Name
		}
	}
	if req.Discount > req.Price {
		return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "Product Discount cannot exceed Product Price")
	}
	newProduct := &Product{
		Name:        req.Name,
		Slug:        utils.GenerateSlug(req.Name),
		Description: req.Description,
		Images:       req.Images,
		Thumbnail:   req.Thumbnail,
		Price:       req.Price,
		Discount:    req.Discount,
		Status:      req.Status,
		Quantity:    req.Quantity,
		Variant:     req.Variant,
		Store:       req.Store,
		Category:    category.Name,
		Subcategory: subcategoryName,
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
	if len(req.Images) != 0 {
		existingProduct.Images = append(existingProduct.Images, req.Images...)
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

	if len(req.Variant) != 0 {
		existingProduct.Variant = append(existingProduct.Variant, req.Variant...)
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
	if err := r.db.Where("user_id", userId).Find(&wishlist).Error; err != nil {
		return nil, err
	}
	return wishlist, nil
}

func (r *repository) RemoveWishListedProduct(ctx context.Context, id uint32) error {
	existingWishlist := &WishListedProduct{}
	err := r.db.Where("id", id).Delete(existingWishlist).Error
	return err
}

func (r *repository) GetProducts(ctx context.Context, store string, limit int, offset int) ([]*Product, error) {
	var products []*Product
	query := r.db
	if store != "" {
		query = query.Where("store_id = ?", store)
	} else {
		query = query.Where("status = ?", true)
	}

	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
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

func (r *repository) AddRecommendedProducts(ctx context.Context, userId uint32, productId uint32) error {
	_, err := r.GetProduct(ctx, productId)
	if err != nil {
		return err
	}
	// existingProduct.Views = append(existingProduct.Views, userId)
	// err = r.db.Model(existingProduct).Update("views", existingProduct.Views).Error
	return err
}

func (r *repository) SearchProducts(ctx context.Context, query string) ([]*Product, error) {
	var products []*Product
	if err := r.db.Where("name ILIKE ? OR category ILIKE ? OR subcategory ILIKE ? OR store ILIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *repository) GetRecommendedProducts(ctx context.Context, query string) ([]*Product, error) {
	var products []*Product
	if err := r.db.Where("category ILIKE ?", "%"+query+"%").Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
