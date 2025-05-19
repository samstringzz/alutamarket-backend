package product

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/Chrisentech/aluta-market-api/database"
	"github.com/Chrisentech/aluta-market-api/errors"
	"github.com/Chrisentech/aluta-market-api/graph/model"
	"github.com/Chrisentech/aluta-market-api/utils"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository() Repository {
	return &repository{
		db: database.GetDB(),
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

	// Update all provided fields
	if req.Name != "" {
		existingProduct.Name = req.Name
		existingProduct.Slug = utils.GenerateSlug(req.Name)
	}
	if req.Description != "" {
		existingProduct.Description = req.Description
	}
	if len(req.Images) > 0 {
		existingProduct.Images = req.Images
	}
	if req.Thumbnail != "" {
		existingProduct.Thumbnail = req.Thumbnail
	}
	if req.Price != 0 {
		existingProduct.Price = req.Price
	}
	if req.Discount != 0 {
		// Validate discount doesn't exceed price
		if req.Discount > existingProduct.Price {
			return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "Product Discount cannot exceed Product Price")
		}
		existingProduct.Discount = req.Discount
	}
	if req.Status != nil {
		existingProduct.Status = *req.Status
	}
	if req.Quantity != 0 {
		existingProduct.Quantity = req.Quantity
	}
	if req.File != "" {
		existingProduct.File = req.File
	}
	if req.Store != "" {
		existingProduct.Store = req.Store
	}
	if req.CategoryID != 0 {
		// Get category to validate and get its name
		category, err := r.GetCategory(ctx, uint32(req.CategoryID))
		if err != nil {
			return nil, fmt.Errorf("failed to get category: %v", err)
		}
		existingProduct.Category = category.Name
	}
	if req.SubCategoryName != "" {
		existingProduct.Subcategory = req.SubCategoryName
	}

	// Update the existing record with all changes
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

func (r *repository) GetProducts(ctx context.Context, storeName string, categorySlug string, limit int, offset int) ([]*Product, int, error) {
	var products []*Product
	var totalCount int64

	// Create base query
	query := r.db.
		Table("products").
		Joins("LEFT JOIN stores ON CAST(products.store AS INTEGER) = stores.id").
		Where("products.deleted_at IS NULL").
		Where("stores.maintenance_mode = ?", false) // Exclude products from stores in maintenance mode

	if storeName != "" {
		// Get store ID from store name
		var storeID uint32
		if err := r.db.Table("stores").Where("name = ?", storeName).Select("id").Scan(&storeID).Error; err != nil {
			return nil, 0, fmt.Errorf("failed to find store: %v", err)
		}

		// Convert storeID to string since products.store is character varying
		storeIDStr := strconv.FormatUint(uint64(storeID), 10)
		query = query.Where("products.store = ?", storeIDStr)
	}

	if categorySlug != "" {
		query = query.
			Joins("JOIN categories ON categories.name = products.category").
			Where("categories.slug = ?", categorySlug)
	} else {
		// For random ordering, use a different approach
		query = query.Where("products.status = ?", true)
		if limit > 0 {
			query = query.Order("random()").Limit(limit)
		}
	}

	// Count total records before pagination
	countQuery := query
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination if not already limited by random selection
	if categorySlug != "" || storeName != "" {
		query = query.Limit(limit).Offset(offset * limit)
	}

	// Execute the final query
	if err := query.Find(&products).Error; err != nil {
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

	// Sanitize and optimize the search query
	formattedQuery := "%" + strings.ToLower(strings.TrimSpace(query)) + "%"

	// Use more efficient query with indexes and exclude products from stores in maintenance mode
	err := r.db.Select("DISTINCT ON (products.id) products.*").
		Table("products").
		Joins("LEFT JOIN categories ON LOWER(categories.name) = LOWER(products.category)").
		Joins("LEFT JOIN stores ON CAST(products.store AS INTEGER) = stores.id"). // Fix: Cast store to integer
		Where("products.deleted_at IS NULL").
		Where("stores.maintenance_mode = ?", false).
		Where("LOWER(products.name) ILIKE ? OR LOWER(products.category) ILIKE ? OR LOWER(COALESCE(categories.slug, '')) ILIKE ?",
			formattedQuery, formattedQuery, formattedQuery).
		Order("products.id, products.created_at DESC").
		Find(&products).Error

	if err != nil {
		log.Printf("Search products error: %v", err)
		return nil, fmt.Errorf("failed to search products: %v", err)
	}

	// Batch initialize slices for better performance
	for _, p := range products {
		if p == nil {
			continue
		}
		p.Images = make([]string, 0, 5) // Preallocate with capacity
		p.Variant = make([]*VariantType, 0, 3)
		p.Views = make([]uint32, 0, 10)
		p.Reviews = make([]Review, 0, 5)
	}

	return products, nil
}

func (r *repository) GetRecommendedProducts(ctx context.Context, query string) ([]*Product, error) {
	var products []*Product
	err := r.db.
		Table("products").
		Joins("LEFT JOIN stores ON products.store = stores.id").
		Where("stores.maintenance_mode = ?", false). // Exclude products from stores in maintenance mode
		Where("category ILIKE ?", "%"+query+"%").
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

// AddReview adds a new review for a product.
func (r *repository) AddReview(ctx context.Context, review *Review) (*Review, error) {
	if err := r.db.Create(review).Error; err != nil {
		return nil, err
	}
	return review, nil
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
			AlwaysAvailable: &p.AlwaysAvailbale,
			Quantity:        p.Quantity,
			File:            file,
			Store:           p.Store,
			Category:        p.Category,
			Subcategory:     p.Subcategory,
			Type:            &p.Type,
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

func (r *repository) GetAllProducts(ctx context.Context) ([]*Product, error) {
	var products []*Product
	err := r.db.
		Table("products").
		Joins("LEFT JOIN stores ON products.store = stores.id").
		Where("stores.maintenance_mode = ?", false). // Exclude products from stores in maintenance mode
		Find(&products).Error
	if err != nil {
		return nil, err
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

	return products, nil
}
