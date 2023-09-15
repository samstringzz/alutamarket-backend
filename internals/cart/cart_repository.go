package cart

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Chrisentech/aluta-market-api/utils"
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
func (r *repository) AddToCart(ctx context.Context, req []*CartItems, user uint32) (*Cart, error) {
	// Check if the user's cart exists and is active
	var cart Cart
	var items []interface{}
	err := r.db.Where("user_id = ? AND active = ?", user, true).First(&cart).Error
	if err != nil {
		for _, item := range req {
			items = append(items, item)
		}

		cartItemsJSON := utils.MarshalJSON(items)
		fmt.Println(cartItemsJSON)

		// Product does not exist in the cart, add it
		// cart.Items = cartItemsJSON
		cart.Total += utils.CalculateTotalCartCost(items) //req.Product.Price
		cart.UserID = user
		cart.Active = true
		fmt.Println(&cart)
		err = r.db.Save(&cart).Error
		if err != nil {
			return nil, err
		}

		return &cart, nil
	}
	cartItemsJSON := utils.MarshalJSON(items)
	fmt.Println(cartItemsJSON)

	for i, item := range req {
		if item.Product.ID == req[i].Product.ID {
			// Product exists, update the quantity
			item.Quantity += item.Quantity
		}
		cart.Total += float64(item.Quantity) * 10 //req.Product.Price
	}

	// Marshal the cart.Items slice to JSON
	// cartItemsJSON := utils.MarshalJSON(cart.Items)

	// Unmarshal the JSON data back into the cart.Items field
	err = utils.UnmarshalJSON(cartItemsJSON, &cart.Items)
	if err != nil {
		return nil, err
	}

	// Save the changes to the database
	err = r.db.Save(&cart).Error
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *repository) GetCart(ctx context.Context, user uint32) (*Cart, error) {
	var cart Cart
	items := utils.MarshalJSON(cart.Items)
	query := r.db.Where("active = true").Where("user_id = ?", user)
	if err := query.First(&cart).Error; err != nil {
		return nil, err
	}
	if err := utils.UnmarshalJSON(items, &cart.Items); err != nil {
		// Handle the error, e.g., log it or return an error response
		fmt.Println("Error unmarshaling JSON:", err)
	}

	return &cart, nil
}

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
