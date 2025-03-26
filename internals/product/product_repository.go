package product

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Chrisentech/aluta-market-api/errors"
	"github.com/Chrisentech/aluta-market-api/graph/model"
	"github.com/Chrisentech/aluta-market-api/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository() Repository {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	// Configure PostgreSQL connection using Supabase connection string
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dbURL,
	}), &gorm.Config{})

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

func (r *repository) GetProduct(ctx context.Context, id, user uint32) (*Product, error) {
	p := &Product{}
	err := r.db.Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}

	// Check if the product is nil
	if p == nil {
		return nil, fmt.Errorf("product not found")
	}

	if user != 0 {
		_, err = r.AddHandledProduct(ctx, user, id, "recently_viewed")
		return nil, err
	} else {
		return p, nil
	}
}

func (r *repository) CreateProduct(ctx context.Context, req *NewProduct) (*Product, error) {
	category, err := r.GetCategory(ctx, uint32(req.CategoryID))
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %v", err)
	}

	if category == nil {
		return nil, fmt.Errorf("category with ID %d not found", req.CategoryID)
	}

	// Verify subcategory exists in the category (check both name and slug)
	subcategoryExists := false
	for _, sub := range category.SubCategories {
		// Compare both name and slug case-insensitively
		if strings.EqualFold(sub.Name, req.SubCategoryName) ||
			strings.EqualFold(sub.Slug, req.SubCategoryName) {
			subcategoryExists = true
			// Use the correct name from the database
			req.SubCategoryName = sub.Name
			break
		}
	}

	if !subcategoryExists {
		return nil, fmt.Errorf("subcategory '%s' not found in category. Available subcategories: %v",
			req.SubCategoryName,
			getSubcategoryNames(category.SubCategories))
	}

	// Ensure images array is not nil
	if req.Images == nil {
		req.Images = []string{}
	}

	// Ensure variant is properly handled
	if req.Variant == nil {
		req.Variant = []*VariantType{}
	}
	if req.Discount > req.Price {
		return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "Product Discount cannot exceed Product Price")
	}
	newProduct := &Product{
		Name:            req.Name,
		Slug:            utils.GenerateSlug(req.Name),
		Description:     req.Description,
		Images:          req.Images,
		Thumbnail:       req.Thumbnail,
		Price:           req.Price,
		Discount:        req.Discount,
		Status:          *req.Status,
		Quantity:        req.Quantity,
		File:            req.File,
		Variant:         req.Variant,
		Store:           req.Store,
		Category:        category.Name,
		Subcategory:     req.SubCategoryName,
		AlwaysAvailbale: req.AlwaysAvailbale,
	}
	if err := r.db.Create(newProduct).Error; err != nil {

		log.Printf("Error creating product: %v", err)
		return nil, err
	}
	return newProduct, nil
}

// Helper function to get subcategory names for error message
func getSubcategoryNames(subs []SubCategory) []string {
	names := make([]string, len(subs))
	for i, sub := range subs {
		names[i] = sub.Name
	}
	return names
}

func (r *repository) UpdateProduct(ctx context.Context, req *NewProduct) (*Product, error) {
	idUint32, err := strconv.ParseUint(req.ID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %v", err)
	}

	// Get existing product
	var existingProduct Product
	if err := r.db.WithContext(ctx).Where("id = ?", idUint32).First(&existingProduct).Error; err != nil {
		return nil, err
	}

	// Update status if provided
	if req.Status != nil {
		existingProduct.Status = *req.Status
	}

	// Update the existing record
	if err := r.db.WithContext(ctx).Model(&existingProduct).Where("id = ?", idUint32).Updates(&existingProduct).Error; err != nil {
		return nil, fmt.Errorf("failed to update product: %v", err)
	}

	return &existingProduct, nil
}

func (r *repository) DeleteProduct(ctx context.Context, id uint32) error {
	existingProduct, err := r.GetProduct(ctx, id, 0)
	if err != nil {
		fmt.Println("Error retrieving product:", err)
		return fmt.Errorf("error retrieving product: %v", err)
	}
	if existingProduct == nil {
		return fmt.Errorf("product not found")
	}
	r.db.Unscoped().Delete(&existingProduct)

	return nil
}

func (r *repository) AddHandledProduct(ctx context.Context, userId, productId uint32, eventType string) (*HandledProduct, error) {
	prd := &HandledProduct{}
	foundProduct, err := r.GetProduct(ctx, productId, 0)
	if err != nil {
		return nil, err
	}

	// Validate the eventType
	if eventType != "recently_viewed" && eventType != "wishlists" && eventType != "savedItems" {
		return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "Type allowed are recently_viewed, wishlists, and savedItems only")
	}

	var count int64
	r.db.Model(prd).Where("product_id = ? AND type = ? AND user_id = ?", productId, eventType, userId).Count(&count)

	// If the product exists and the eventType is recently_viewed, return the existing product
	if count > 0 {
		if eventType == "recently_viewed" {
			r.db.Where("product_id = ? AND type = ? AND user_id = ?", productId, eventType, userId).First(prd)
			return prd, nil
		} else {
			return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "Product already exists for this type")
		}
	}

	// If the product doesn't exist, create a new one
	prd.Product = foundProduct
	prd.UserID = userId
	prd.Type = eventType
	err = r.db.Create(prd).Error
	if err != nil {
		return nil, err
	}
	return prd, nil
}

func (r *repository) GetHandledProducts(ctx context.Context, userId uint32, eventType string) ([]*HandledProduct, error) {
	var prds []*HandledProduct
	if err := r.db.Where("user_id = ? AND type = ? ", userId, eventType).Find(&prds).Error; err != nil {
		return nil, err
	}
	// Print the fetched products for debugging
	// for _, prd := range prds {
	// 	log.Printf("Fetched product: %+v\n", prd)
	// }
	return prds, nil
}

func (r *repository) RemoveHandledProduct(ctx context.Context, id uint32, eventType string) error {
	handlePrd := &HandledProduct{}
	err := r.db.Unscoped().Where("id=? AND type=?", id, eventType).Delete(handlePrd).Error
	return err
}

func (r *repository) RemoveWishListedProduct(ctx context.Context, id uint32) error {
	existingWishlist := &HandledProduct{}
	err := r.db.Where("id=? ", id).Delete(existingWishlist).Error
	return err
}

func (r *repository) GetProducts(ctx context.Context, store string, categorySlug string, limit int, offset int) ([]*Product, int, error) {
	var products []*Product
	var totalCount int64

	query := r.db
	if store != "" {
		query = query.Where("store = ?", store)
	}
	if categorySlug != "" {
		query = query.Joins("JOIN categories ON categories.name = products.category").
			Where("categories.slug = ?", categorySlug)
	} else {
		query = query.Where("status = ?", true).
			Order("RANDOM()")
	}

	// Count total records
	if err := query.Model(&Product{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated products
	if err := query.Limit(limit).Offset(offset * limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	// Initialize JSON fields
	for _, p := range products {
		if p.Images == nil {
			p.Images = []string{}
		}
		if p.Variant == nil {
			p.Variant = []*VariantType{}
		}
		if p.Views == nil {
			p.Views = []uint32{}
		}
		if p.Reviews == nil {
			p.Reviews = []Review{}
		}
	}

	return products, int(totalCount), nil
}

func (r *repository) GetCategories(ctx context.Context) ([]*Category, error) {
	var categories []*Category
	if err := r.db.Preload("SubCategories").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("error retrieving categories: %v", err)
	}
	return categories, nil
}

func (r *repository) GetCategory(ctx context.Context, id uint32) (*Category, error) {
	p := Category{}
	err := r.db.Preload("SubCategories").Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// isValidQuery validates the query string for allowed characters
func isValidQuery(query string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9- ]+$", query)
	return match
}

func filterReviewsBySeller(reviews []*Review, sellerId string) []*Review {
	var filteredReviews []*Review
	for _, review := range reviews {
		if review.Seller == sellerId {
			filteredReviews = append(filteredReviews, review)
		}
	}
	return filteredReviews
}

func (r *repository) SearchProducts(ctx context.Context, query string) ([]*Product, error) {
	var products []*Product

	// Ensure the query string is properly formatted
	formattedQuery := "%" + query + "%"

	// Perform the search operation
	err := r.db.Select("products.*").
		Table("products").
		Joins("JOIN categories ON categories.name = products.category").
		Where("(categories.slug ILIKE ? OR products.name ILIKE ?) AND products.deleted_at IS NULL",
			formattedQuery,
			formattedQuery).
		Find(&products).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search products: %v", err)
	}

	// Initialize empty slices if nil
	for _, p := range products {
		if p.Images == nil {
			p.Images = []string{}
		}
		if p.Variant == nil {
			p.Variant = []*VariantType{}
		}
		if p.Views == nil {
			p.Views = []uint32{}
		}
		if p.Reviews == nil {
			p.Reviews = []Review{}
		}
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

// AddReview adds a new review for a product.
func (r *repository) AddReview(ctx context.Context, input *Review) (*Review, error) {
	input.ID = utils.GenerateUUID() // Generate unique ID for the review

	// Save the review directly to the database
	err := r.db.Create(input).Error
	if err != nil {
		return nil, fmt.Errorf("failed to add review: %v", err)
	}

	return input, nil
}

// GetProductReviews retrieves all reviews for a specific product, or all reviews for a seller's products.
func (r *repository) GetProductReviews(ctx context.Context, productId uint32, sellerId string) ([]*Review, error) {
	var reviews []*Review

	// If we are fetching reviews for a specific product
	if productId > 0 {
		err := r.db.Where("product_id = ?", productId).Find(&reviews).Error
		if err != nil {
			return nil, fmt.Errorf("failed to fetch product reviews: %v", err)
		}
		return reviews, nil
	}

	// If we are fetching reviews for all products of a seller
	// In GetProductReviews function
	if sellerId != "" {
		var products []*Product
		// Add empty string for categorySlug parameter
		products, _, err := r.GetProducts(ctx, sellerId, "", 10000, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch seller's products: %v", err)
		}

		productIDs := make([]uint32, len(products))
		for i, prd := range products {
			productIDs[i] = prd.ID
		}

		err = r.db.Where("product_id IN ?", productIDs).Find(&reviews).Error
		if err != nil {
			return nil, fmt.Errorf("failed to fetch reviews for seller's products: %v", err)
		}
		return reviews, nil
	}

	return reviews, nil
}

func (r *repository) AddSavedForLater(ctx context.Context, userId, productId uint32) (*HandledProduct, error) {
	savedForLater := &HandledProduct{}
	foundProduct, err := r.GetProduct(ctx, productId, 0)
	if err != nil {
		return nil, err
	}
	var count int64
	r.db.Model(savedForLater).Where("user_id =?", userId).Count(&count)
	if count > 0 {
		fmt.Printf("The Total no of User savedForLater is%v\n", count)
		return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "Product already in savedForLater")
	}
	savedForLater.Product = foundProduct
	savedForLater.UserID = userId
	savedForLater.Type = "savedForLater"
	err = r.db.Create(savedForLater).Error
	if err != nil {
		return nil, err
	}
	return savedForLater, nil
}

func (r *repository) GetSavedForLaterProducts(ctx context.Context, userId uint32) ([]*HandledProduct, error) {
	var savedForLater []*HandledProduct
	if err := r.db.Where("user_id = ? AND type = ? ", userId, "savedForLater").Find(&savedForLater).Error; err != nil {
		return nil, err
	}
	return savedForLater, nil
}

func (r *repository) Products(ctx context.Context, store *string, categorySlug *string, limit *int, offset *int) (*model.ProductPaginationData, error) {
	storeValue := ""
	if store != nil {
		storeValue = *store
	}

	categorySlugValue := ""
	if categorySlug != nil {
		categorySlugValue = *categorySlug
	}

	limitValue := 10
	if limit != nil {
		limitValue = *limit
	}

	offsetValue := 0
	if offset != nil {
		offsetValue = *offset
	}

	products, total, err := r.GetProducts(ctx, storeValue, categorySlugValue, limitValue, offsetValue)
	if err != nil {
		return nil, err
	}

	modelProducts := make([]*model.Product, len(products))
	for i, p := range products {
		var file *string
		if p.File != "" {
			file = &p.File
		}

		modelProducts[i] = &model.Product{
			ID:              int(p.ID),
			Name:            p.Name,
			Slug:            p.Slug,
			Description:     p.Description,
			Image:           p.Images, // Use the Images field directly
			Thumbnail:       p.Thumbnail,
			Price:           p.Price,
			Discount:        p.Discount,
			Status:          p.Status,
			AlwaysAvailable: p.AlwaysAvailbale,
			Quantity:        p.Quantity,
			File:            file, // Handle nil case for File field
			Store:           p.Store,
			Category:        p.Category,
			Subcategory:     p.Subcategory,
			Type:            p.Type,
		}
	}

	return &model.ProductPaginationData{
		Data:  modelProducts,
		Total: total,
	}, nil
}

func (r *repository) GetHandledProductsWithDetails(ctx context.Context, userID uint32, typeArg string) ([]*HandledProduct, error) {
	var handledProducts []*HandledProduct

	// Join handled_products with products table
	result := r.db.WithContext(ctx).
		Table("handled_products hp").
		Select("hp.*, p.*").
		Joins("LEFT JOIN products p ON hp.product_id = p.id").
		Where("hp.user_id = ? AND hp.type = ? AND hp.deleted_at IS NULL", userID, typeArg).
		Find(&handledProducts)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch handled products: %v", result.Error)
	}

	return handledProducts, nil
}
