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

func (r *repository) GetProduct(ctx context.Context, id, user uint32) (*Product, error) {
	p := Product{}
	err := r.db.Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	// if user != 0 {
	// 	err = r.AddRecentlyViewedProducts(ctx, user, id)
	// 	return nil, err
	// }
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
		Images:      req.Images,
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
	existingProduct, err := r.GetProduct(ctx, req.ID, 0)
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
	existingProduct, err := r.GetProduct(ctx, id, 0)
	if err != nil {
		return err
	}
	err = r.db.Delete(existingProduct).Error
	return err

}
// func (r *repository) AddWishListedProduct(ctx context.Context, userId, productId uint32) (*HandledProduct, error) {
// 	wishlist := &HandledProduct{}
// 	foundProduct, err := r.GetProduct(ctx, productId, 0)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var count int64
// 	r.db.Model(wishlist).Where("user_id =?", userId).Count(&count)
// 	if count > 0 {
// 		fmt.Printf("The Total no of User wishlist is%v\n", count)
// 		return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "Product already in wishlist")
// 	}
// 	wishlist.Product = foundProduct
// 	wishlist.UserID = userId
// 	wishlist.Type   = "wishlist"
// 	err = r.db.Create(wishlist).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return wishlist, nil
// }

func (r *repository) AddHandledProduct(ctx context.Context, userId, productId uint32,eventType string) (*HandledProduct, error) {
	prd := &HandledProduct{}
	foundProduct, err := r.GetProduct(ctx, productId, 0)
	if err != nil {
		return nil, err
	}
	var count int64
	r.db.Model(prd).Where("user_id =? AND type?=", userId,eventType).Count(&count)
	if count > 0 {
		fmt.Printf("The Total no of User %v\n is%v\n", eventType,count)
		return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "Product already exist for this type ")
	}
	
	prd.Product = foundProduct
	prd.UserID = userId
	prd.Type   = eventType
	err = r.db.Create(prd).Error
	if err != nil {
		return nil, err
	}
	return prd, nil
}

func (r *repository) GetHandledProducts(ctx context.Context, userId uint32,eventType string) ([]*HandledProduct, error) {
	var prds []*HandledProduct
	if err := r.db.Where("user_id = ? AND type = ? ", userId,eventType).Find(&prds).Error; err != nil {
		return nil, err
	}
	return prds, nil
}
// func (r *repository) GetWishListedProducts(ctx context.Context, userId uint32) ([]*HandledProduct, error) {
// 	var wishlist []*HandledProduct
// 	if err := r.db.Where("user_id = ? AND type = ? ", userId,"wishlist").Find(&wishlist).Error; err != nil {
// 		return nil, err
// 	}
// 	return wishlist, nil
// }

func (r *repository) RemoveHandledProduct(ctx context.Context, id uint32,eventType string) error {
	existingWishlist := &HandledProduct{}
	err := r.db.Where("id=? AND type=?", id,eventType).Delete(existingWishlist).Error
	return err
}

func (r *repository) RemoveWishListedProduct(ctx context.Context, id uint32) error {
	existingWishlist := &HandledProduct{}
	err := r.db.Where("id=? ", id).Delete(existingWishlist).Error
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
// func (r *repository) AddRecentlyViewedProducts(ctx context.Context, userId, productId uint32) error {
// 	product, err := r.GetProduct(ctx, productId, userId)

// 	if err != nil {
// 		return err
// 	}

// 	foundView := false

// 	// Update the Views of the product clicked
// 	for _, view := range product.Views {
// 		if view == userId {
// 			foundView = true
// 			break
// 		}
// 	}

// 	if !foundView {
// 		product.Views = append(product.Views, userId)
// 		err = r.db.Save(product).Error
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	// Update the user's recentlyViewed
// 	foundUser, err := user.NewRepository().GetUser(ctx, strconv.FormatUint(uint64(userId), 10)) // Replace with your GetUser implementation

// 	if err != nil {
// 		return err
// 	}

// 	// Check if the product already exists in the RecentlyViewed
// 	var productIndex uint32 = 0
// 		for i, viewedProduct := range foundUser.RecentlyViewed {
// 		if viewedProduct.ID == productId {
// 			productIndex = uint32(i)
// 			break
// 		}
// 	}
	

// 	if productIndex >0 {
// 		// Move the product to the first item in the array
// 		viewedProduct := foundUser.RecentlyViewed[productIndex]
// 		foundUser.RecentlyViewed = append(foundUser.RecentlyViewed[:productIndex], foundUser.RecentlyViewed[productIndex+1:]...)
// 		foundUser.RecentlyViewed = append(foundUser.RecentlyViewed, viewedProduct)
// 	} else {
// 	fmt.Println("The length of the recently viewed product is: ", len(foundUser.RecentlyViewed))

// 		// If the product is not already in RecentlyViewed, add it
// 		recentlyViewedProduct := product
// 		foundUser.RecentlyViewed = append(foundUser.RecentlyViewed, recentlyViewedProduct)
// 	}
// 	// Ensure the RecentlyViewed list is not more than 8 items
// 	if len(foundUser.RecentlyViewed) > 8 {
// 		foundUser.RecentlyViewed = foundUser.RecentlyViewed[:8]
// 	}

// 	// Update the user in the database (assuming a SaveUser function)
// 	err = r.db.Save(foundUser).Error // Replace with your SaveUser implementation

// 	return nil
// }

// func (r *repository) GetRecentlyViewedProducts(ctx context.Context, userId uint32) ([]*Product, error) {
// 	// Update the user's recentlyViewed
// 	foundUser, err := user.NewRepository().GetUser(ctx, strconv.FormatUint(uint64(userId), 10)) // Replace with your GetUser implementation
// 	if err != nil {
// 		return nil, err
// 	}
// 	viewedProducts := foundUser.RecentlyViewed
// 	var products []*Product

// 	for _, item := range viewedProducts {
// 		product, _ := r.GetProduct(ctx, item.ID, 0)
// 		products = append(products, product)
// 	}
// 	return products, nil
// }

func (r *repository) SearchProducts(ctx context.Context, query string) ([]*Product, error) {
	var result []*Product
	if err := r.db.Where("name ILIKE ? OR category ILIKE ? OR subcategory ILIKE ? OR store ILIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) GetRecommendedProducts(ctx context.Context, query string) ([]*Product, error) {
	var products []*Product
	if err := r.db.Where("category ILIKE ?", "%"+query+"%").Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *repository) AddReview(ctx context.Context, input *Review) (*Review, error) {
	product, err := r.GetProduct(ctx, input.ProductID, 0)
	if err != nil {
		return nil, err
	}
	product.Reviews = append(product.Reviews, input)
	err = r.db.Save(product).Error
	if err != nil {
		return nil, err
	}
	return input, nil
}

func (r *repository) GetReviews(ctx context.Context, productId uint32) ([]*Review, error) {
	product, err := r.GetProduct(ctx, productId, 0)
	if err != nil {
		return nil, err
	}
	return product.Reviews, nil
}

func (r *repository) AddSavedForLater(ctx context.Context, userId,productId uint32)(*HandledProduct,error){
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
	savedForLater.Type   = "savedForLater"
	err = r.db.Create(savedForLater).Error
	if err != nil {
		return nil, err
	}
	return savedForLater, nil
}

func (r *repository) GetSavedForLaterProducts(ctx context.Context, userId uint32) ([]*HandledProduct, error) {
	var savedForLater []*HandledProduct
	if err := r.db.Where("user_id = ? AND type = ? ", userId,"savedForLater").Find(&savedForLater).Error; err != nil {
		return nil, err
	}
	return savedForLater, nil
}
