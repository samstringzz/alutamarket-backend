package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

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

func NewRepository() *repository {
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

type SquadWebhookPayload struct {
	Event string `json:"event"`
	Data  struct {
		ID          string `json:"id"`          // Unique transaction ID
		Amount      int64  `json:"amount"`      // Amount of the transaction
		Currency    string `json:"currency"`    // Currency of the transaction
		Status      string `json:"status"`      // Status of the transaction (e.g., "completed")
		Reference   string `json:"reference"`   // Reference code for the transaction
		Description string `json:"description"` // Description or narration of the transaction
		CreatedAt   string `json:"created_at"`  // Timestamp when the transaction was created
		UpdatedAt   string `json:"updated_at"`  // Timestamp when the transaction was last updated
		Customer    struct {
			ID          string `json:"id"`           // Customer ID
			Name        string `json:"name"`         // Customer's name
			Email       string `json:"email"`        // Customer's email
			PhoneNumber string `json:"phone_number"` // Customer's phone number
		} `json:"customer"`
		PaymentMethod struct {
			Type string `json:"type"` // Type of payment method (e.g., "card")
			Card struct {
				Number   string `json:"number"`    // Card number (may be masked)
				Last4    string `json:"last4"`     // Last 4 digits of the card number
				ExpMonth string `json:"exp_month"` // Expiration month of the card
				ExpYear  string `json:"exp_year"`  // Expiration year of the card
				CardType string `json:"card_type"` // Type of the card (e.g., "visa")
			} `json:"card"`
		} `json:"payment_method"`
		Metadata struct {
			CustomFields []struct {
				DisplayName string `json:"display_name"` // Custom field display name
				Value       string `json:"value"`        // Custom field value
			} `json:"custom_fields"`
		} `json:"metadata"`
	} `json:"data"`
}

type PSWebhookPayload struct {
	Event string `json:"event"`
	Data  struct {
		ID              int64  `json:"id"`
		Domain          string `json:"domain"`
		Status          string `json:"status"`
		Reference       string `json:"reference"`
		Amount          int64  `json:"amount"`
		Message         string `json:"message"`
		GatewayResponse string `json:"gateway_response"`
		PaidAt          string `json:"paid_at"`
		CreatedAt       string `json:"created_at"`
		Channel         string `json:"channel"`
		Currency        string `json:"currency"`
		IPAddress       string `json:"ip_address"`
		Metadata        struct {
			CustomFields []struct {
				DisplayName string `json:"display_name"`
				Value       string `json:"value"`
			} `json:"custom_fields"`
		} `json:"metadata"`
		Customer struct {
			ID           int64  `json:"id"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Email        string `json:"email"`
			CustomerCode string `json:"customer_code"`
		} `json:"customer"`
		Authorization struct {
			AuthorizationCode string `json:"authorization_code"`
			Bin               string `json:"bin"`
			Last4             string `json:"last4"`
			ExpMonth          string `json:"exp_month"`
			ExpYear           string `json:"exp_year"`
			CardType          string `json:"card_type"`
			Bank              string `json:"bank"`
			CountryCode       string `json:"country_code"`
		} `json:"authorization"`
	} `json:"data"`
}

func ConvertTrackedToStoreProduct(tp store.TrackedProduct) *store.StoreProduct {
	return &store.StoreProduct{
		Name:      tp.Name,
		Price:     tp.Price,
		Quantity:  tp.Quantity,
		Thumbnail: tp.Thumbnail,
		ID:        tp.ID,
	}
}

type FWWebhookPayload struct {
	Event string `json:"event"`
	Data  struct {
		ID                int64   `json:"id"`
		TxRef             string  `json:"tx_ref"`
		FlwRef            string  `json:"flw_ref"`
		DeviceFingerprint string  `json:"device_fingerprint"`
		Amount            float64 `json:"amount"`
		Currency          string  `json:"currency"`
		ChargedAmount     float64 `json:"charged_amount"`
		AppFee            float64 `json:"app_fee"`
		MerchantFee       float64 `json:"merchant_fee"`
		ProcessorResponse string  `json:"processor_response"`
		AuthModel         string  `json:"auth_model"`
		IP                string  `json:"ip"`
		Narration         string  `json:"narration"`
		Status            string  `json:"status"`
		PaymentType       string  `json:"payment_type"`
		CreatedAt         string  `json:"created_at"`
		AccountID         int64   `json:"account_id"`
		Customer          struct {
			ID          int64  `json:"id"`
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
			Email       string `json:"email"`
			CreatedAt   string `json:"created_at"`
		} `json:"customer"`
		Entity struct {
			AccountNumber string `json:"account_number"`
			AccountBank   string `json:"account_bank"`
		} `json:"entity"`
	} `json:"data"`
}

func VerifyFWWebhookSignature(req *http.Request, secretHash string, body []byte) bool {
	signature := req.Header.Get("verif-hash")
	if signature == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secretHash))
	mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}
func (repo *repository) FWWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Close the request body to avoid resource leaks
	defer r.Body.Close()

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Verify the webhook signature
	secretHash := os.Getenv("FLW_SECRET_HASH") // Replace with your actual Flutterwave secret hash
	if !VerifyFWWebhookSignature(r, secretHash, body) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Unmarshal the JSON data into a struct
	var webhookPayload FWWebhookPayload
	if err := json.Unmarshal(body, &webhookPayload); err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
	// fmt.Printf("Response: %+v\n", webhookPayload.Data)
	fmt.Printf("Response: %s\n", string(body))

	// Acknowledge the webhook request early
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Flutterwave Webhook received successfully")

	// Access the event field
	event := webhookPayload.Event

	//Access the data field
	data := webhookPayload.Data

	// Process the event
	switch event {

	case "charge.completed":
		fmt.Printf("fwdata: %+v\n", data)
		order := &store.Order{}
		sellerOrder := &store.StoreOrder{}
		seller := &user.User{}
		buyer := &user.User{}

		var productsWithFiles []store.TrackedProduct

		// Handle charge completion logic here

		//fecth  order
		err := repo.db.Model(order).Where("uuid", data.TxRef).Error
		if err != nil {
			http.Error(w, "Failed to find order", http.StatusNotFound)
			return
		}

		uniqueStores := utils.RemoveDuplicates(order.StoresID)
		for _, storeName := range uniqueStores {
			// Fetch the seller corresponding to the store name
			err := repo.db.Model(&seller).Where("name = ?", storeName).First(&seller).Error
			if err != nil {
				http.Error(w, "Failed to find seller", http.StatusNotFound)
				return
			}

			if data.Status == "successful" {
				// Filter products that have associated files
				for _, product := range order.Products {
					if product.File != nil {
						productsWithFiles = append(productsWithFiles, product)
					}
				}

				// Process each product with a file
				for _, product := range productsWithFiles {
					download := &store.Downloads{
						Thumbnail: product.Thumbnail,
						Price:     product.Price,
						Name:      product.Name,
						Discount:  int(product.Discount),
						UUID:      order.UUID,
						File:      *product.File,          // Dereference the file pointer safely
						Users:     []string{order.UserID}, // Initialize with the current user
					}

					// Save the download record to the database
					if err := repo.db.Create(download).Error; err != nil {
						http.Error(w, "Failed to save download", http.StatusInternalServerError)
						return
					}
					storeProduct := ConvertTrackedToStoreProduct(product)
					sellerOrder.Products = append(sellerOrder.Products, storeProduct)
				}
				sellerOrder.Customer = order.Customer
				sellerOrder.Status = "pending"
				sellerOrder.UUID = order.UUID
				sellerOrder.Active = true
				// Save the seller's order
				if err := repo.db.Save(sellerOrder).Error; err != nil {
					http.Error(w, "Failed to update seller order", http.StatusInternalServerError)
					return
				}

				// Send email notification to the seller
				if seller.Email != "" {
					to := []string{seller.Email}
					contents := map[string]string{
						"seller_name":     seller.Fullname,
						"order_id":        data.TxRef,
						"products_length": strconv.Itoa(len(order.Products)),
						"customer_name":   buyer.Fullname,
						"customer_phone":  buyer.Phone,
					}

					templateID := "991b93a9-4661-452c-ba53-da31fdddf8f2"
					utils.SendEmail(templateID, "New Order Alert🎉", to, contents)
				}
			}
		}

		err = repo.db.Model(buyer).Where("id", order.UserID).Error
		if err != nil {
			http.Error(w, "Failed to find user", http.StatusNotFound)
			return
		}

		// Pay Delivery fee
		err = user.PayFund(float32(order.Fee), seller.Email, "3002290305", "50211")
		if err != nil {
			http.Error(w, "Failed to pay delivery", http.StatusNotFound)
			return
		}

		order.TransStatus = data.Status
		order.Status = "pending"
		order.PaymentMethod = data.PaymentType
		order.PaymentGateway = "flutterwave"
		// Email credentials

		if err := repo.db.Save(order).Error; err != nil {
			http.Error(w, "Failed to update  order", http.StatusInternalServerError)
			return
		}

		fmt.Printf("Received charge completed event: %s for email: %s, Amount: %d\n", event, webhookPayload.Data.Customer.Email, int64(webhookPayload.Data.Amount))

	default:
		// Optionally handle other events or log them
		fmt.Println("Received unhandled event:", event)
	}

	// Log the received body for debugging or further processing
	fmt.Println("Received Webhook Body:", string(body))
}

func VerifyPSWebhookSignature(req *http.Request, secret string, body []byte) bool {
	signature := req.Header.Get("x-paystack-signature")
	if signature == "" {
		return false
	}

	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

func (repo *repository) PaystackWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Close the request body to avoid resource leaks
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Verify the webhook signature
	secret := os.Getenv("PAYSTACK_SECRET_KEY") // Replace with your actual Paystack secret key
	if !VerifyPSWebhookSignature(r, secret, body) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Unmarshal the JSON data into a struct
	var webhookPayload PSWebhookPayload
	if err := json.Unmarshal(body, &webhookPayload); err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	fmt.Printf("Response: %s\n", string(body))

	// Acknowledge the webhook request early
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Paystack Webhook received successfully")

	// Access the event field
	event := webhookPayload.Event
	data := webhookPayload.Data

	order := &store.Order{}
	sellerOrder := &store.StoreOrder{}
	seller := &user.User{}
	buyer := &user.User{}
	myStore := &store.Store{}

	var productsWithFiles []store.TrackedProduct
	// Process the event
	switch event {
	case "charge.success":
		fmt.Printf("psdata: %+v\n", data)

		// Handle charge success logic here

		//fecth  order
		err := repo.db.Model(order).Where("uuid", data.Reference).Error
		if err != nil {
			http.Error(w, "Failed to find order", http.StatusNotFound)
			return
		}

		// Fetch Stores
		uniqueStores := utils.RemoveDuplicates(order.StoresID)
		for _, storeName := range uniqueStores {
			// Fetch the seller corresponding to the store name
			err := repo.db.Model(&myStore).Where("name = ?", storeName).First(&myStore).Error
			if err != nil {
				http.Error(w, "Failed to find seller", http.StatusNotFound)
				return
			}

			if data.Status == "success" {
				// Filter products that have associated files
				for _, product := range order.Products {
					if product.File != nil {
						productsWithFiles = append(productsWithFiles, product)
					}
				}

				// Process each product with a file
				for _, product := range productsWithFiles {
					download := &store.Downloads{
						Thumbnail: product.Thumbnail,
						Price:     product.Price,
						Name:      product.Name,
						Discount:  int(product.Discount),
						UUID:      order.UUID,
						File:      *product.File,          // Dereference the file pointer safely
						Users:     []string{order.UserID}, // Initialize with the current user
					}

					// Save the download record to the database
					if err := repo.db.Create(download).Error; err != nil {
						http.Error(w, "Failed to save download", http.StatusInternalServerError)
						return
					}
					storeProduct := ConvertTrackedToStoreProduct(product)
					sellerOrder.Products = append(sellerOrder.Products, storeProduct)
				}
				sellerOrder.Customer = order.Customer
				sellerOrder.Status = "pending"
				sellerOrder.UUID = order.UUID
				sellerOrder.Active = true
				// Save the seller's order
				if err := repo.db.Save(sellerOrder).Error; err != nil {
					http.Error(w, "Failed to update seller order", http.StatusInternalServerError)
					return
				}

				err := repo.db.Model(seller).Where("id", myStore.UserID).Error
				if err != nil {
					http.Error(w, "Failed to find seller", http.StatusNotFound)
					return
				}

				// Send email notification to the seller
				if seller.Email != "" {
					to := []string{seller.Email}
					contents := map[string]string{
						"seller_name":     seller.Fullname,
						"order_id":        data.Reference,
						"products_length": strconv.Itoa(len(order.Products)),
						"customer_name":   buyer.Fullname,
						"customer_phone":  buyer.Phone,
					}

					templateID := "991b93a9-4661-452c-ba53-da31fdddf8f2"
					utils.SendEmail(templateID, "New Order Alert🎉", to, contents)
				}
			}
		}
		err = repo.db.Model(buyer).Where("id", order.UserID).Error
		if err != nil {
			http.Error(w, "Failed to find user", http.StatusNotFound)
			return
		}

		// Pay Delivery fee
		err = user.PayFund(float32(order.Fee), seller.Email, "3002290305", "50211")
		if err != nil {
			http.Error(w, "Failed to pay delivery", http.StatusNotFound)
			return
		}
		order.TransStatus = data.Status
		order.Status = "pending"
		order.PaymentMethod = data.Channel
		order.PaymentGateway = "paystack"
		sellerOrder.Active = true

		if err := repo.db.Save(order).Error; err != nil {
			http.Error(w, "Failed to update seller order", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Received charge success event: %s for email: %s, Amount: %d\n", event, webhookPayload.Data.Customer.Email, webhookPayload.Data.Amount)

	case "charge.failed":
		// Handle charge failure logic here
		// Handle charge success logic here
		err := repo.db.Model(order).Where("trt_ref = ? ", data.Reference).Error
		if err != nil {
			http.Error(w, "Failed to find order", http.StatusNotFound)
			return
		}
		order.TransStatus = data.Status
		order.Status = "pending"
		order.PaymentMethod = data.Channel
		order.PaymentGateway = "paystack"
		if err := repo.db.Save(order).Error; err != nil {
			http.Error(w, "Failed to save order", http.StatusInternalServerError)
			return
		}

		fmt.Printf("Received charge failed event: %s for email: %s\n", event, webhookPayload.Data.Customer.Email)

	default:
		// Optionally handle other events or log them
		fmt.Println("Received unhandled event:", event)
	}

	// Log the received body for debugging or further processing
	fmt.Println("Received Webhook Body:", string(body))
}

func (repo *repository) SquadWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Close the request body to avoid resource leaks
	defer r.Body.Close()

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// You can implement signature verification here if Squad provides such a mechanism

	// Unmarshal the JSON data into a struct
	var webhookPayload SquadWebhookPayload
	if err := json.Unmarshal(body, &webhookPayload); err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	// Acknowledge the webhook request early
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Squad Webhook received successfully")

	// Access the event field
	event := webhookPayload.Event
	data := webhookPayload.Data

	buyerOrder := &store.Order{}
	sellerOrder := &store.StoreOrder{}
	seller := &user.User{}
	buyer := &user.User{}
	myStore := &store.Store{}
	var productsWithFiles []store.TrackedProduct

	// Process the event
	switch event {

	case "payment.success":
		fmt.Printf("sqdata: %+v\n", data)
		// Handle charge success logic here
		err := repo.db.Model(buyerOrder).Where("trt_ref = ? ", data.Reference).Error
		if err != nil {
			http.Error(w, "Failed to find order", http.StatusNotFound)
			return
		}
		// Fetch Store
		err = repo.db.Model(myStore).Where("id", sellerOrder.StoreID).Error
		if err != nil {
			http.Error(w, "Failed to find user", http.StatusNotFound)
			return
		}

		err = repo.db.Model(seller).Where("id", myStore.UserID).Error
		if err != nil {
			http.Error(w, "Failed to find user", http.StatusNotFound)
			return
		}

		err = repo.db.Model(buyer).Where("id", myStore.UserID).Error
		if err != nil {
			http.Error(w, "Failed to find user", http.StatusNotFound)
			return
		}

		// Pay Delivery fee
		err = user.PayFund(float32(buyerOrder.Fee), seller.Email, "3002290305", "50211")
		if err != nil {
			http.Error(w, "Failed to pay delivery", http.StatusNotFound)
			return
		}
		buyerOrder.TransStatus = data.Status
		buyerOrder.Status = "pending"
		buyerOrder.PaymentMethod = data.PaymentMethod.Type
		buyerOrder.PaymentGateway = "squad"
		sellerOrder.Active = true
		for _, product := range buyerOrder.Products {
			if product.File != nil {
				productsWithFiles = append(productsWithFiles, product)
			}
		}
		for _, product := range productsWithFiles {
			// Check if the file is not nil before dereferencing
			if product.File != nil {
				download := &store.Downloads{
					Thumbnail: product.Thumbnail,
					Price:     product.Price,
					Name:      product.Name,
					Discount:  int(product.Discount),
					UUID:      buyerOrder.UUID,
					File:      *product.File, // Safely dereferencing product.File
				}

				// Initialize Users slice if it's nil
				if download.Users == nil {
					download.Users = make([]string, 0)
				}

				// Append the user ID
				download.Users = append(download.Users, buyerOrder.UserID)

				// Save download to database
				repo.db.Create(download)
			}
		}
		to := []string{seller.Email}
		contents := map[string]string{
			"seller_name":     seller.Fullname,
			"order_id":        sellerOrder.UUID,
			"products_length": strconv.Itoa(len(sellerOrder.Products)),
			"customer_name":   buyer.Fullname,
			"customer_phone":  buyer.Phone,
		}

		templateID := "991b93a9-4661-452c-ba53-da31fdddf8f2"
		utils.SendEmail(templateID, "New Order Alert🎉", to, contents)

		if err := repo.db.Save(buyerOrder).Error; err != nil {
			http.Error(w, "Failed to save order", http.StatusInternalServerError)
			return
		}
		if err := repo.db.Save(sellerOrder).Error; err != nil {
			http.Error(w, "Failed to update seller order", http.StatusInternalServerError)
			return
		}
		// Handle successful payment logic here
		fmt.Printf("Payment success: %s for email: %s, Amount: %d\n", webhookPayload.Data.Reference, webhookPayload.Data.Customer.Email, webhookPayload.Data.Amount)

	case "payment.failed":
		// Handle failed payment logic here
		buyerOrder := &store.Order{}
		// Handle charge success logic here
		err := repo.db.Model(buyerOrder).Where("trt_ref = ? ", data.Reference).Error
		if err != nil {
			http.Error(w, "Failed to find order", http.StatusNotFound)
			return
		}
		buyerOrder.TransStatus = data.Status
		buyerOrder.Status = "pending"
		buyerOrder.PaymentMethod = data.PaymentMethod.Type
		buyerOrder.PaymentGateway = "squad"
		if err := repo.db.Save(buyerOrder).Error; err != nil {
			http.Error(w, "Failed to save order", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Payment failed: %s for email: %s\n", webhookPayload.Data.Reference, webhookPayload.Data.Customer.Email)

	default:
		// Optionally handle other events or log them
		fmt.Println("Received unhandled event:", event)
	}

	// Log the received body for debugging or further processing
	fmt.Println("Received Webhook Body:", string(body))
}
