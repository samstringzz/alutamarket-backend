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
	"time"

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
	w.WriteHeader(http.StatusOK)

	go func() {
		// Unmarshal the JSON data into a struct
		var webhookPayload FWWebhookPayload
		if err := json.Unmarshal(body, &webhookPayload); err != nil {
			http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
			return
		}
		// fmt.Printf("Response: %+v\n", webhookPayload.Data)
		fmt.Printf("Response: %s\n", string(body))

		// Acknowledge the webhook request early

		fmt.Fprint(w, "Flutterwave Webhook received successfully")

		// Access the event field
		event := webhookPayload.Event

		//Access the data field
		data := webhookPayload.Data

		order := &store.Order{}
		sellerOrder := &store.StoreOrder{}
		buyer := &user.User{}
		var productsWithFiles []store.TrackedProduct

		// Process the event
		switch event {

		case "charge.completed":

			//fecth  order
			err := repo.db.Model(order).Where("uuid=?", data.TxRef).First(order).Error
			if err != nil {
				http.Error(w, "Failed to find order", http.StatusNotFound)
				return
			}

			// Update product statuses to "pending"
			for i := range order.Products {
				order.Products[i].Status = "pending"
			}

			// Save the updated products back to the database
			// if err := repo.db.Save(&order.Products).Error; err != nil {
			// 	fmt.Println("Failed to update product statuses:", err)
			// 	return
			// }

			err = repo.db.Model(buyer).Where("id = ?", order.UserID).First(buyer).Error
			if err != nil {
				fmt.Println("Failed to find buyer:", err)
				return
			}

			// Fetch Stores
			uniqueStores := utils.RemoveDuplicatesStringArray([]string(order.StoresID))

			for _, storeName := range uniqueStores {
				myStore := &store.Store{} // Ensure a clean struct
				seller := &user.User{}
				// Fetch the store corresponding to the store name
				if storeName != nil {
					err := repo.db.Model(&myStore).Where("name = ?", *storeName).First(&myStore).Error
					if err != nil {
						fmt.Println("Failed to find store:", err)
					}
				}

				// Fetch seller
				err = repo.db.Model(&seller).Where("id = ?", myStore.UserID).First(&seller).Error
				if err != nil {
					fmt.Println("Failed to find seller:", err)
					continue
				}

				if data.Status == "successful" {
					// Filter products that have associated files
					for _, product := range order.Products {
						if product.File != nil && *product.File != "" {
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
							ID:        utils.GenerateUUID(),
							File:      *product.File,          // Dereference the file pointer safely
							Users:     []string{order.UserID}, // Initialize with the current user
						}

						// Save the download record to the database
						if err := repo.db.Create(download).Error; err != nil {
							http.Error(w, "Failed to save download", http.StatusInternalServerError)
							return
						}
						// storeProduct := ConvertTrackedToStoreProduct(product)
						// sellerOrder.Products = append(sellerOrder.Products, storeProduct)
					}

					sellerOrder.Products = filterProductsByStore(order.Products, *storeName)
					sellerOrder.StoreID = strconv.Itoa(int(myStore.ID))
					sellerOrder.Customer = order.Customer
					sellerOrder.Status = "pending"
					sellerOrder.UUID = order.UUID
					sellerOrder.Active = true
					sellerOrder.CreatedAt = time.Now() // Ensure createdAt is explicitly se

					myStore.Orders = append(myStore.Orders, sellerOrder)
					// Save the seller's store
					if err := repo.db.Save(myStore).Error; err != nil {
						fmt.Println("Failed to update seller order:", err)
						continue
					}
					// fmt.Printf("Updated Seller Order: %+v\n", sellerOrder)

					// Send email notification to the seller
					if seller.Email != "" {
						to := []string{seller.Email}
						contents := map[string]string{
							"seller_name":     seller.Fullname,
							"order_id":        order.UUID,
							"products_length": strconv.Itoa(len(order.Products)),
							"customer_name":   buyer.Fullname,
							"customer_phone":  buyer.Phone,
						}

						templateID := "991b93a9-4661-452c-ba53-da31fdddf8f2"
						utils.SendEmail(templateID, "New Order Alert🎉", to, contents)
						fmt.Printf("Email sent to seller: %s\n", seller.Email)

					}
				}

			}

			// Pay delivery fee
			user.PayFund(float32(order.Fee), "3002290305", "50211")
			// if err != nil {
			// 	fmt.Println("Failed to pay delivery fee:", err)
			// 	return
			// }
			// fmt.Println("Delivery fee paid successfully")

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

		fmt.Println("Webhook processing completed")

	}()
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
	var secret string
	if os.Getenv("ENVIRONMENT") == "development" {
		secret = os.Getenv("PAYSTACK_SECRET_KEY_DEV")
	} else {
		secret = os.Getenv("PAYSTACK_SECRET_KEY")
	}
	if !VerifyPSWebhookSignature(r, secret, body) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Send the HTTP 200 response immediately
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Paystack Webhook received")
	go func() {
		// Continue processing the webhook in a goroutine

		// Unmarshal the JSON data into a struct
		var webhookPayload PSWebhookPayload
		if err := json.Unmarshal(body, &webhookPayload); err != nil {
			fmt.Println("Failed to parse JSON data:", err)
			return
		}

		// Log the raw payload for debugging
		// fmt.Printf("Webhook Payload: %s\n", string(body))

		// Access the event field
		event := webhookPayload.Event
		data := webhookPayload.Data

		order := &store.Order{}
		sellerOrder := &store.StoreOrder{}
		buyer := &user.User{}

		var productsWithFiles []store.TrackedProduct

		// Process the event
		switch event {
		case "charge.success":
			// fmt.Printf("Charge success event received: %+v\n", data)

			// Fetch order
			err := repo.db.Model(order).Where("uuid = ?", data.Reference).First(order).Error
			if err != nil {
				fmt.Println("Failed to find order:", err)
				return
			}

			// Update product statuses to "pending"
			for i := range order.Products {
				order.Products[i].Status = "pending"
			}

			// Save the updated products back to the database
			if err := repo.db.Save(&order).Error; err != nil {
				fmt.Println("Failed to update product statuses:", err)
				return
			}
			fmt.Println("Product statuses updated to 'pending'")
			// fmt.Printf("Fetched Order: %+v\n", order)

			err = repo.db.Model(buyer).Where("id = ?", order.UserID).First(buyer).Error
			if err != nil {
				fmt.Println("Failed to find buyer:", err)
				return
			}
			// fmt.Printf("Fetched Buyer: %+v\n", buyer)

			// Fetch Stores
			uniqueStores := utils.RemoveDuplicatesStringArray([]string(order.StoresID))

			for _, storeName := range uniqueStores {
				myStore := &store.Store{} // Ensure a clean struct
				seller := &user.User{}
				// Fetch the store details
				if storeName != nil {
					err := repo.db.Model(&myStore).Where("name = ?", *storeName).First(&myStore).Error
					if err != nil {
						fmt.Println("Failed to find store:", err)
					}
				}
				// fmt.Printf("Fetched Store: %+v\n", myStore)

				// Fetch seller
				err = repo.db.Model(&seller).Where("id = ?", myStore.UserID).First(&seller).Error
				if err != nil {
					fmt.Println("Failed to find seller:", err)
					continue
				}
				// fmt.Printf("Fetched Seller: %+v\n", seller)

				if data.Status == "success" {
					// Filter products that have associated files
					for _, product := range order.Products {
						if product.File != nil && *product.File != "" {
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
							ID:        utils.GenerateUUID(),
							File:      *product.File,          // Dereference the file pointer safely
							Users:     []string{order.UserID}, // Initialize with the current user
						}

						// Save the download record to the database
						if err := repo.db.Create(download).Error; err != nil {
							http.Error(w, "Failed to save download", http.StatusInternalServerError)
							return
						}
						// storeProduct := ConvertTrackedToStoreProduct(product)
						// sellerOrder.Products = append(sellerOrder.Products, storeProduct)
					}

					sellerOrder.Products = filterProductsByStore(order.Products, *storeName)
					sellerOrder.StoreID = strconv.Itoa(int(myStore.ID))
					sellerOrder.Customer = order.Customer
					sellerOrder.Status = "pending"
					sellerOrder.UUID = order.UUID
					sellerOrder.Active = true
					sellerOrder.CreatedAt = time.Now() // Ensure createdAt is explicitly se

					myStore.Orders = append(myStore.Orders, sellerOrder)

					// Save the seller's store
					if err := repo.db.Save(myStore).Error; err != nil {
						fmt.Println("Failed to update seller order:", err)
						continue
					}
					// fmt.Printf("Updated Seller Order: %+v\n", sellerOrder)

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
						fmt.Printf("Email sent to seller: %s\n", seller.Email)

					}
				}
			}

			// Pay delivery fee
			user.PayFund(float32(order.Fee), "3002290305", "50211")
			// if err != nil {
			// 	fmt.Println("Failed to pay delivery fee:", err)
			// 	return
			// }
			// fmt.Println("Delivery fee paid successfully")

			// Update order
			order.TransStatus = data.Status
			order.Status = "pending"
			order.PaymentMethod = data.Channel
			order.PaymentGateway = "paystack"
			if err := repo.db.Save(order).Error; err != nil {
				fmt.Println("Failed to update order:", err)
				return
			}
			// fmt.Printf("Updated Order: %+v\n", order)

			fmt.Printf("Charge success handled successfully: Event: %s, Email: %s, Amount: %d\n",
				event, webhookPayload.Data.Customer.Email, webhookPayload.Data.Amount)

		case "charge.failed":
			// Handle charge failure
			fmt.Printf("Charge failed event received: %+v\n", data)
			err := repo.db.Model(order).Where("uuid = ?", data.Reference).First(order).Error
			if err != nil {
				fmt.Println("Failed to find order:", err)
				return
			}
			order.TransStatus = data.Status
			order.Status = "failed"
			order.PaymentMethod = data.Channel
			order.PaymentGateway = "paystack"
			if err := repo.db.Save(order).Error; err != nil {
				fmt.Println("Failed to save order:", err)
				return
			}
			fmt.Printf("Updated Order for Charge Failed: %+v\n", order)

		default:
			// Log unhandled events
			fmt.Println("Received unhandled event:", event)
		}

		fmt.Println("Webhook processing completed")
	}()
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
	w.WriteHeader(http.StatusOK)

	// Unmarshal the JSON data into a struct
	go func() {
		var webhookPayload SquadWebhookPayload
		if err := json.Unmarshal(body, &webhookPayload); err != nil {
			http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
			return
		}

		// Acknowledge the webhook request early
		fmt.Fprint(w, "Squad Webhook received successfully")

		// Access the event field
		event := webhookPayload.Event
		data := webhookPayload.Data

		order := &store.Order{}
		sellerOrder := &store.StoreOrder{}
		buyer := &user.User{}

		var productsWithFiles []store.TrackedProduct

		// Process the event
		switch event {

		case "payment.success":
			// fmt.Printf("sqdata: %+v\n", data)
			// Handle charge success logic here
			// Fetch order
			err := repo.db.Model(order).Where("uuid = ?", data.Reference).First(order).Error
			if err != nil {
				fmt.Println("Failed to find order:", err)
				return
			}
			// Fetch Store
			err = repo.db.Model(buyer).Where("id = ?", order.UserID).First(buyer).Error
			if err != nil {
				fmt.Println("Failed to find buyer:", err)
				return
			}
			// Fetch Stores
			uniqueStores := utils.RemoveDuplicatesStringArray([]string(order.StoresID))

			for _, storeName := range uniqueStores {
				myStore := &store.Store{} // Ensure a clean struct
				seller := &user.User{}
				// Fetch the store details
				if storeName != nil {
					err := repo.db.Model(&myStore).Where("name = ?", *storeName).First(&myStore).Error
					if err != nil {
						fmt.Println("Failed to find store:", err)
					}
				}
				// fmt.Printf("Fetched Store: %+v\n", myStore)

				// Fetch seller
				err = repo.db.Model(&seller).Where("id = ?", myStore.UserID).First(&seller).Error
				if err != nil {
					fmt.Println("Failed to find seller:", err)
					continue
				}
				// fmt.Printf("Fetched Seller: %+v\n", seller)

				if data.Status == "successful" {
					// Filter products that have associated files
					for _, product := range order.Products {
						if product.File != nil && *product.File != "" {
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
							ID:        utils.GenerateUUID(),
							File:      *product.File,          // Dereference the file pointer safely
							Users:     []string{order.UserID}, // Initialize with the current user
						}

						// Save the download record to the database
						if err := repo.db.Create(download).Error; err != nil {
							http.Error(w, "Failed to save download", http.StatusInternalServerError)
							return
						}
						// storeProduct := ConvertTrackedToStoreProduct(product)
						// sellerOrder.Products = append(sellerOrder.Products, storeProduct)
					}

					sellerOrder.Products = filterProductsByStore(order.Products, *storeName)
					sellerOrder.StoreID = strconv.Itoa(int(myStore.ID))
					sellerOrder.Customer = order.Customer
					sellerOrder.Status = "pending"
					sellerOrder.UUID = order.UUID
					sellerOrder.Active = true
					sellerOrder.CreatedAt = time.Now() // Ensure createdAt is explicitly se

					myStore.Orders = append(myStore.Orders, sellerOrder)
					// Save the seller's store
					if err := repo.db.Save(myStore).Error; err != nil {
						fmt.Println("Failed to update seller order:", err)
						continue
					}

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
						fmt.Printf("Email sent to seller: %s\n", seller.Email)

					}
				}
			}

			// Pay delivery fee
			user.PayFund(float32(order.Fee), "3002290305", "50211")
			// if err != nil {
			// 	fmt.Println("Failed to pay delivery fee:", err)
			// 	return
			// }
			// Update order
			order.TransStatus = data.Status
			order.Status = "pending"
			order.PaymentMethod = data.PaymentMethod.Type
			order.PaymentGateway = "paystack"
			if err := repo.db.Save(order).Error; err != nil {
				fmt.Println("Failed to update order:", err)
				return
			}
			// fmt.Printf("Updated Order: %+v\n", order)

			fmt.Printf("Charge success handled successfully: Event: %s, Email: %s, Amount: %d\n",
				event, webhookPayload.Data.Customer.Email, webhookPayload.Data.Amount)

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
	}()
}

func filterProductsByStore(products []store.TrackedProduct, storeName string) []*store.StoreProduct {
	var filteredProducts []*store.StoreProduct

	for i, product := range products {
		product.Quantity = products[i].Quantity
		storeProduct := ConvertTrackedToStoreProduct(product)
		if product.Store == storeName {
			filteredProducts = append(filteredProducts, storeProduct)
		}
	}

	return filteredProducts
}
