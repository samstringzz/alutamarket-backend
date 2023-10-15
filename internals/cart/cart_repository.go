package cart

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"encoding/json"
	"log"
    "net/http"
    "io"
	"os"

	"github.com/Chrisentech/aluta-market-api/errors"
	"github.com/Chrisentech/aluta-market-api/internals/product"
	"github.com/Chrisentech/aluta-market-api/internals/store"
	"github.com/Chrisentech/aluta-market-api/internals/user"
	"github.com/Chrisentech/aluta-market-api/utils"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}
type WebhookPayload struct {
    Event string `json:"event"`
    Data  string `json:"data"`
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

func (r *repository) ModifyCart(ctx context.Context, req *CartItems, user uint32) (*Cart, error) {
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
func (r *repository) RemoveAllCart(ctx context.Context, id uint32) (*Cart, error) {
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

func (r *repository) InitiatePayment(ctx context.Context, input Order) (string, error) {
	UUID := utils.GenerateUUID()
	redirectUrl := "http://yemi.com/redirect"
	cart := &Cart{}
	userID, _ := strconv.ParseUint(input.UserID, 10, 32)

	err := r.db.Model(cart).Where("user_id = ? AND active =?", uint32(userID),true).First(cart).Error
	// fmt.Println("Na here we dey")
	if err != nil {
		return "", errors.NewAppError(http.StatusNotFound, "NOT FOUND", "Cart not found")
	}
	customer, _ := user.NewRepository().GetUser(ctx, input.UserID)
	requestData := map[string]interface{}{
		"tx_ref":       UUID,
		"amount":       input.Fee + cart.Total,
		"currency":     "NGN",
		"redirect_url": redirectUrl,
		"meta": map[string]interface{}{
			"consumer_id":  23,
			"consumer_mac": "92a3-912ba-1192a",
		},
		"customer": map[string]interface{}{
			"email":       customer.Email,
			"phonenumber": customer.Phone,
			"name":        customer.Fullname,
		},
		"customizations": map[string]interface{}{
			"title": "Aluta Market Checkout",
			"logo":  "http://www.piedpiper.com/app/themes/joystick-v27/images/logo.png",
		},
	}
	// fmt.Println("And not here")

	// Serialize the request data to JSON
	requestDataBytes, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	// Create an HTTP client
	client := &http.Client{}

	// Create the HTTP request
	req, err := http.NewRequest("POST", "https://api.flutterwave.com/v3/payments", bytes.NewBuffer(requestDataBytes))
	if err != nil {
		return "", err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("FLW_SECRET_KEY"))

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	fmt.Println(response)
	var paymentLink string
	responseData, ok := response["data"].(map[string]interface{})
	if ok {
		linkValue, linkExists := responseData["link"].(string)
		if linkExists {
			paymentLink = linkValue
		}
	}
	newOrder := &store.Order{}
	newOrder.Amount = cart.Total + input.Fee
	newOrder.UserID = strconv.FormatUint(uint64(cart.UserID), 10)
	newOrder.UUID = UUID
	newOrder.Status = "pending"
	for _, items := range cart.Items {
		newOrder.StoresID = append(newOrder.StoresID, items.Product.Store)
	}
	if input.PaymentGateway == "flutterwave" {
		newOrder.PaymentGateway = "flutterwave"
		newOrder.Status = ""
	} else if input.PaymentGateway == "paystack" {
		newOrder.PaymentGateway = "paystack"
		// return("paystack")
	}
	cart.Active = false
	r.db.Save(&cart)
	err = r.db.Model(&store.Order{}).Save(newOrder).Error
	if err != nil {
		return "", err
	}
	return paymentLink, nil
}

func (r *repository) MakePayment(ctx context.Context,w http.ResponseWriter, req *http.Request){
	body, err :=io.ReadAll(req.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }
    // Unmarshal the JSON data into a struct
    var webhookPayload WebhookPayload
    
   if err := json.Unmarshal(body, &webhookPayload); err != nil {
        http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
        return
    }

    // Access the event field
    event := webhookPayload.Event

    if event =="charge.completed"{
        return
    }

    // Close the request body to avoid resource leaks
    defer req.Body.Close()

    // Process the request body (e.g., decode JSON or parse form data)
    
    fmt.Fprint(w, "Webhook received successfully")
    fmt.Println("Received Webhook Body:", string(body))

	// Check if transaction is successful/failed,by rehitting the verify transaction api
	// If Successful, credit a store wallet,their quota

     w.WriteHeader(http.StatusOK)
}
