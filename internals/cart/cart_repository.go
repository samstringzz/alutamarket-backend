package cart

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Chrisentech/aluta-market-api/errors"
	"github.com/Chrisentech/aluta-market-api/internals/product"
	"github.com/Chrisentech/aluta-market-api/internals/store"
	"github.com/Chrisentech/aluta-market-api/internals/user"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}
type WebhookPayload struct {
	Event string `json:"event"`
	Data  struct {
		TransactionID  string  `json:"transaction_id"`
		Amount         float64 `json:"amount"`
		Currency       string  `json:"currency"`
		TransactionRef string  `json:"tx_ref"`
	} `json:"data"`
}

type VerifyResponse struct {
	Data struct {
		Status   string  `json:"status"`
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	} `json:"data"`
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
	var err2 error
	if req.Product.ID == 0 && req.Product.Name == "" {
		return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "both id and name cannot be empty")
	} else if req.Product.ID == 0 {
		err2 = r.db.Model(&prd).
			Where("name = ?", req.Product.Name).
			First(&prd).Error
	} else if req.Product.Name == "" {
		err2 = r.db.Model(&prd).
			Where("id = ?", req.Product.ID).
			First(&prd).Error
	} else {
		err2 = r.db.Model(&prd).
			Where("name = ? OR id = ?", req.Product.Name, req.Product.ID).
			First(&prd).Error
	}
	newQuantity := prd.Quantity + (-req.Quantity)

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

func (r *repository) GetCart(ctx context.Context, filter uint32) (*Cart, error) {
	var cart *Cart
	query := r.db.Where("active = true").Where("user_id = ?", filter)
	if err := query.First(&cart).Error; err != nil {
		err = r.db.Where("active = true").Where("id = ?", filter).First(&cart).Error
		if err != nil {
			return &Cart{}, nil
		}
	}

	return cart, nil
}

func (r *repository) RemoveAllCart(ctx context.Context, id uint32) error {
	// Find the cart with the specified ID that is active
	var cart Cart
	query := r.db.Where("active = true").Where("id = ?", id)
	if err := query.First(&cart).Error; err != nil {
		return err
	}

	// Iterate over the items in the cart and update the product quantities
	for _, item := range cart.Items {
		var product product.Product
		if err := r.db.Where("id = ?", item.Product.ID).First(&product).Error; err != nil {
			return err
		}
		product.Quantity += item.Quantity
		if err := r.db.Save(&product).Error; err != nil {
			return err
		}
	}

	// Set the 'active' field to false and update the cart
	cart.Active = false
	if err := r.db.Save(&cart).Error; err != nil {
		return err
	}

	return nil
}

func processPaymentGateway(gateway string, UUID string, Fee float64, cartTotal float64, redirectUrl string, customer *user.User) (string, error) {
	var paymentLink string
	var err error
	var requestData map[string]interface{}
	var req *http.Request
	//  cart.Total + input.Amount
	// Set the payment gateway key based on the environment
	var gatewayKey string
	if os.Getenv("ENVIRONMENT") == "development" {
		switch gateway {
		case "flutterwave":
			gatewayKey = os.Getenv("FLW_SECRET_KEY_DEV")
		case "paystack":
			gatewayKey = os.Getenv("PAYSTACK_SECRET_KEY_DEV")
		case "squad":
			gatewayKey = os.Getenv("SQUAD_SECRET_KEY_DEV")
		}
	} else {
		switch gateway {
		case "flutterwave":
			gatewayKey = os.Getenv("FLW_SECRET_KEY")
		case "paystack":
			gatewayKey = os.Getenv("PAYSTACK_SECRET_KEY")
		case "squad":
			gatewayKey = os.Getenv("SQUAD_SECRET_KEY")
		}
	}

	// Prepare request data based on the payment gateway
	switch gateway {
	case "flutterwave":
		requestData = map[string]interface{}{
			"tx_ref":       UUID,
			"amount":       Fee + cartTotal,
			"currency":     "NGN",
			"redirect_url": redirectUrl,
			"meta": map[string]interface{}{
				"consumer_id":  customer.ID,
				"consumer_mac": "92a3-912ba-1192a",
			},
			"customer": map[string]interface{}{
				"email":       customer.Email,
				"phonenumber": customer.Phone,
				"name":        customer.Fullname,
			},
			"customizations": map[string]interface{}{
				"title": "Aluta Market Checkout",
				"logo":  "https://res.cloudinary.com/folajimidev/image/upload/v1697737213/logo_xesoiu.png",
			},
		}

		req, err = http.NewRequest("POST", "https://api.flutterwave.com/v3/payments", nil)

	case "paystack":
		requestData = map[string]interface{}{
			"email":        customer.Email,
			"amount":       (Fee + cartTotal) * 100, // Amount should be in kobo
			"currency":     "NGN",
			"callback_url": redirectUrl,
			"metadata": map[string]interface{}{
				"consumer_id":  customer.ID,
				"consumer_mac": "92a3-912ba-1192a",
			},
			"customizations": map[string]interface{}{
				"title": "Aluta Market Checkout",
				"logo":  "https://res.cloudinary.com/folajimidev/image/upload/v1697737213/logo_xesoiu.png",
			},
		}

		req, err = http.NewRequest("POST", "https://api.paystack.co/transaction/initialize", nil)

	case "squad":
		requestData = map[string]interface{}{
			"currency":        "NGN",
			"initiate_type":   "inline",
			"transaction_ref": UUID,
			"callback_url":    redirectUrl,
			"amount":          (Fee + cartTotal) * 100,
			"email":           customer.Email,
			"pass_charge":     true,
			"customer_name":   customer.Fullname,
		}
		// Serialize the request data to JSON
		squadRequestDataBytes, err := json.Marshal(requestData)
		if err != nil {
			return "", err
		}
		req, err = http.NewRequest("POST", "https://sandbox-api-d.squadco.com/transaction/initiate", bytes.NewBuffer(squadRequestDataBytes))
		if err != nil {
			return "", err
		}

	default:
		return "", errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "unsupported payment gateway")
	}

	// Check if there was an error creating the request
	if err != nil {
		return "", err
	}

	// Serialize the request data to JSON
	requestDataBytes, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+gatewayKey)
	req.Body = io.NopCloser(bytes.NewBuffer(requestDataBytes))

	// Make the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Log the response for debugging
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Printf("%s Response: %s\n", gateway, string(respBody))

	// Parse the response
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", err
	}

	responseData, ok := response["data"].(map[string]interface{})
	if ok {
		switch gateway {
		case "flutterwave":
			paymentLink, _ = responseData["link"].(string)
		case "paystack":
			paymentLink, _ = responseData["authorization_url"].(string)

		case "squad":
			paymentLink, _ = responseData["checkout_url"].(string)
		}
	}

	return paymentLink, nil
}

func (r *repository) InitiatePayment(ctx context.Context, input Order) (string, error) {
	UUID := input.UUID
	redirectUrl := "https://www.thealutamarket.com"

	errChan := make(chan error, 2)
	cartChan := make(chan *Cart)
	customerChan := make(chan *user.User)
	userID, _ := strconv.ParseUint(input.UserID, 10, 32)

	// Fetch cart concurrently
	go func() {
		cart := &Cart{}
		err := r.db.Model(cart).Where("user_id = ? AND active =?", uint32(userID), true).First(cart).Error
		if err != nil {
			errChan <- errors.NewAppError(http.StatusNotFound, "NOT FOUND", "Cart not found")
		} else {
			cartChan <- cart
		}
	}()

	// Fetch customer details concurrently
	go func() {
		customer, err := user.NewRepository().GetUser(ctx, input.UserID)
		if err != nil {
			errChan <- err
		} else {
			customerChan <- customer
		}
	}()

	var cart *Cart
	var customer *user.User
	for i := 0; i < 2; i++ {
		select {
		case err := <-errChan:
			return "", err
		case cart = <-cartChan:
		case customer = <-customerChan:
		}
	}

	// Payment gateway processing (this is I/O-bound, consider running it concurrently)
	paymentLinkChan := make(chan string)
	paymentErrChan := make(chan error)

	go func() {
		paymentLink, err := processPaymentGateway(input.PaymentGateway, UUID, input.Amount, cart.Total, redirectUrl, customer)
		if err != nil {
			paymentErrChan <- err
		} else {
			paymentLinkChan <- paymentLink
		}
	}()

	// Handle concurrent payment processing
	var paymentLink string
	select {
	case paymentLink = <-paymentLinkChan:
	case err := <-paymentErrChan:
		return "", err
	}

	newOrder := &store.Order{}
	newOrder.Amount = cart.Total + input.Amount
	newOrder.UserID = strconv.FormatUint(uint64(cart.UserID), 10)
	newOrder.UUID = UUID
	newOrder.TransStatus = "not processed"
	newOrder.Status = "canceled"
	newOrder.Fee = input.Amount
	newOrder.CartID = cart.ID
	newOrder.CreatedAt = time.Now().UTC()
	newOrder.UpdatedAt = time.Now().UTC()
	newOrder.PaymentGateway = input.PaymentGateway
	newOrder.DeliveryDetails = store.DeliveryDetails{
		Method:  customer.PaymentDetails.Info,
		Address: customer.PaymentDetails.Address,
		Fee:     input.Amount,
	}

	var products []store.TrackedProduct
	for _, item := range cart.Items {
		product := store.TrackedProduct{
			ID:        item.Product.ID,
			Name:      item.Product.Name,
			Thumbnail: item.Product.Thumbnail,
			Price:     item.Product.Price,
			Discount:  item.Product.Discount,
			Quantity:  item.Product.Quantity,
			Status:    "open",
			CreatedAt: time.Now(), // Set current time for CreatedAt
			UpdatedAt: time.Now(), // Set current time for UpdatedAt
		}

		products = append(products, product)
	}
	newOrder.Products = products
	newOrder.Status = "pending"
	//Create Seller Order Logic
	var storePrd []*store.StoreProduct
	for _, item := range cart.Items {
		product := &store.StoreProduct{
			ID:        item.Product.ID,
			Name:      item.Product.Name,
			Thumbnail: item.Product.Thumbnail,
			Price:     item.Product.Price,
			Quantity:  item.Product.Quantity,
		}

		storePrd = append(storePrd, product)
	}
	req := &store.StoreOrder{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UUID:      UUID,
		Products:  storePrd,
		Active:    false,
		Customer: store.Customer{
			ID:      uint32(customer.ID),
			Name:    customer.PaymentDetails.Name,
			Phone:   customer.PaymentDetails.Phone,
			Address: customer.PaymentDetails.Address,
			Info:    customer.PaymentDetails.Info,
		},
	}
	// Mark the cart as inactive and process the order concurrently
	if paymentLink != "" {
		cart.Active = false
		cart.Total += input.Amount

		errChan = make(chan error, 3)

		// Save the cart concurrently
		go func() {
			if err := r.db.Save(&cart).Error; err != nil {
				errChan <- err
			} else {
				errChan <- nil
			}
		}()

		// Create the order concurrently
		go func() {
			_, err := store.NewRepository().CreateOrder(ctx, req)
			if err != nil {
				errChan <- err
			} else {
				errChan <- nil
			}
		}()

		// Save the order in the database concurrently
		go func() {
			if err := r.db.Model(&store.Order{}).Save(newOrder).Error; err != nil {
				errChan <- err
			} else {
				errChan <- nil
			}
		}()

		// Wait for all goroutines to finish and check for errors
		for i := 0; i < 3; i++ {
			if err := <-errChan; err != nil {
				return "", err
			}
		}
	}

	return paymentLink, nil
}

func (r *repository) MakePayment(ctx context.Context, w http.ResponseWriter, req *http.Request) {

	body, err := io.ReadAll(req.Body)
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
	var (
		flwSecretKey     = os.Getenv("FLW_SECRET_KEY")
		transactionID    = webhookPayload.Data.TransactionID
		expectedAmount   = webhookPayload.Data.Amount
		expectedCurrency = "NGN"
	)

	// Access the event field
	event := webhookPayload.Event
	//For Flutterwave Charge payment
	if event == "charge.completed" {
		// Check if transaction is successful/failed,by rehitting the verify transaction api
		client := &http.Client{}
		url := "https://api.flutterwave.com/v3/transactions/" + transactionID + "/verify"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating HTTP request:", err)
			return
		}

		req.Header.Add("Authorization", "Bearer "+flwSecretKey)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending HTTP request:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading HTTP response:", err)
			return
		}

		var verifyResponse VerifyResponse
		err = json.NewDecoder(strings.NewReader(string(body))).Decode(&verifyResponse)
		if err != nil {
			fmt.Println("Error decoding JSON response:", err)
			return
		}

		if verifyResponse.Data.Status == "successful" && verifyResponse.Data.Amount == expectedAmount && verifyResponse.Data.Currency == expectedCurrency {
			// Success! Confirm the customer's payment
			order := &store.Order{}
			err = r.db.Where("uuid = ?", webhookPayload.Data.TransactionRef).First(order).Error
			if err != nil {
				fmt.Println("Error Getting order: ", err)
				return
			}
			cart, _ := r.GetCart(ctx, order.CartID)
			for _, item := range cart.Items {
				priceDifference := item.Product.Price - order.Fee
				result, _ := store.NewRepository().GetStoreByName(ctx, item.Product.Store)
				//Credit individual Store from the particular transaction

				result.Wallet += priceDifference

				r.db.Save(result)
			}
			fmt.Println("Payment was successful!")
		} else {
			// Inform the customer their payment was unsuccessful by mail
			fmt.Println("Payment was unsuccessful.")
		}
	}

	// For Flutterwave Transfer Event
	// Close the request body to avoid resource leaks
	defer req.Body.Close()

	// Process the request body (e.g., decode JSON or parse form data)

	fmt.Fprint(w, "Webhook received successfully")
	fmt.Println("Received Webhook Body:", string(body))

	// If Successful, credit a store wallet,their quota

	w.WriteHeader(http.StatusOK)
}
