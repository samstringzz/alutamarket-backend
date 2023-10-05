package cart

import (
	"context"
	// "fmt"
	"log"
	"net/http"
	"os"

	"github.com/Chrisentech/aluta-market-api/errors"
	"github.com/Chrisentech/aluta-market-api/internals/product"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository() Repository {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dbURI := os.Getenv("DB_URI")

	// Initialize the database
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &repository{
		db: db,
	}
}
func calculateTotalCartCost(data []*CartItems) float64 {
	var total float64

	for _, item := range data {
		quantity := item.Quantity
		total += float64(quantity) * (item.Product.Price - item.Product.Discount)
	}
	return total
}

// func (r *repository) AddToCart(ctx context.Context, req *CartItems, user uint32) (*Cart, error) {
// 	prd := &product.Product{}
// 	err2 := r.db.Model(prd).Where("id = ?", req.Product.ID).First(prd).Error
// 	newQuantity := prd.Quantity - req.Quantity

// 	if err2 != nil {
// 		return nil, errors.NewAppError(http.StatusNotFound, "NOT FOUND", "Product not found")
// 	}
// 	if req.Quantity > prd.Quantity {
// 		return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "Product Quantity Exceeded")
// 	}
// 	r.db.Model(prd).Update("quantity", newQuantity)
// 	var cart *Cart
// 	err := r.db.Where("user_id = ? AND active = ?", user, true).First(&cart).Error

// 	if err == nil {
// 		req.Product = prd
// 		// Check if the product is already in the cart
// 		found := false
// 		for _, item := range cart.Items {
// 			//Verify that product being deducted is not above product quantity

// 			if req.Product.ID == item.Product.ID {
// 				if req.Quantity+item.Quantity < 0 {
// 				r.db.Model(prd).Update("quantity", prd.Quantity+req.Quantity)
// 				return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "Product Quantity Exceeded")
// 			}
// 				item.Quantity += req.Quantity
// 				found = true
// 				break
// 			}
// 		}

// 		if !found {
// 			cart.Items = append(cart.Items, req)
// 		}

// 		cart.Total = calculateTotalCartCost(cart.Items)
// 		err = r.db.Save(&cart).Error
// 		if err != nil {
// 			return nil, err
// 		}
// 		return cart, nil
// 	} else {
// 		req.Product.Quantity = newQuantity
// 		req.Product = prd
// 		cart.Items = append(cart.Items, req)
// 		cart.UserID = user
// 		cart.Active = true
// 		cart.Total += float64(req.Quantity) * (prd.Price - prd.Discount)
// 		err = r.db.Save(&cart).Error
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

//		return cart, nil
//	}
func (r *repository) AddToCart(ctx context.Context, req *CartItems, user uint32) (*Cart, error) {
	prd := &product.Product{}
	err2 := r.db.Model(prd).Where("id = ?", req.Product.ID).First(prd).Error
	newQuantity := prd.Quantity - req.Quantity

	if err2 != nil {
		return nil, errors.NewAppError(http.StatusNotFound, "NOT FOUND", "Product not found")
	}
	if req.Quantity > prd.Quantity {
		return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "Product Quantity Exceeded")
	}
	r.db.Model(prd).Update("quantity", newQuantity)
	var cart *Cart
	err := r.db.Where("user_id = ? AND active = ?", user, true).First(&cart).Error

	if err == nil {
		req.Product = prd
		// Check if the product is already in the cart
		found := false
		for i, item := range cart.Items {
			if req.Product.ID == item.Product.ID {
				if req.Quantity+item.Quantity == 0 {
					// Remove the item from the cart when quantity becomes zero or negative
					cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
					r.db.Model(prd).Update("quantity", prd.Quantity+item.Quantity)
				} else if req.Quantity+item.Quantity < 0 {
					r.db.Model(prd).Update("quantity", prd.Quantity+req.Quantity)
					return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "Product Quantity Exceeded")
				} else {
					item.Quantity += req.Quantity
				}
				found = true
				break
			}
		}

		if !found && req.Quantity > 0 {
			// Add the product to the cart only if the quantity is positive
			cart.Items = append(cart.Items, req)
		}

		cart.Total = calculateTotalCartCost(cart.Items)
		err = r.db.Save(&cart).Error
		if err != nil {
			return nil, err
		}
		return cart, nil
	} else {
		req.Product.Quantity = newQuantity
		req.Product = prd
		if req.Quantity > 0 {
			// Add the product to the cart only if the quantity is positive
			cart.Items = append(cart.Items, req)
		}
		cart.UserID = user
		cart.Active = true
		cart.Total += float64(req.Quantity) * (prd.Price - prd.Discount)
		err = r.db.Save(&cart).Error
		if err != nil {
			return nil, err
		}
	}

	return cart, nil
}

func (r *repository) GetCart(ctx context.Context, user uint32) (*Cart, error) {
	var cart *Cart
	query := r.db.Where("active = true").Where("user_id = ?", user)
	if err := query.First(&cart).Error; err != nil {
		return &Cart{}, nil
	}

	return cart, nil
}

// This is not needed
func (r *repository) RemoveFromCart(ctx context.Context, id uint32) (*Cart, error) {
	// Find the cart with the specified ID that is active
	var cart Cart
	query := r.db.Where("active = true").Where("id = ?", id)
	if err := query.First(&cart).Error; err != nil {
		return nil, err
	}

	// Set the 'active' field to false and update the cart
	cart.Active = false
	if err := r.db.Save(&cart).Error; err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *repository) MakePayment(ctx context.Context, req Order) (*Order, error) {
	cart := &Cart{}
	err := r.db.Model(cart).Where("id ?", req.CartID).First(cart).Error
	if err != nil {
		return nil, errors.NewAppError(http.StatusNotFound, "NOT FOUND", "Cart not found")
	}
	// totalCharge:= cart.Total + req.Fee
	if req.PaymentGateway == "flutterwave" {
		panic("flutter")
	} else if req.PaymentGateway == "paystack" {
		panic("paystack")
	}

	panic("nil-not yet implemented")
}
