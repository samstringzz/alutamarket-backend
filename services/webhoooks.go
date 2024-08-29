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

	"github.com/Chrisentech/aluta-market-api/internals/store"
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

		buyerOrder := &store.Order{}
		// Handle charge completion logic here
		err := repo.db.Model(buyerOrder).Where("trt_ref", data.TxRef).Error
		if err != nil {
			http.Error(w, "Failed to find order", http.StatusNotFound)
			return
		}
		buyerOrder.TransStatus = data.Status
		buyerOrder.Status = "pending"
		buyerOrder.PaymentMethod = data.PaymentType
		buyerOrder.PaymentGateway = "flutterwave"
		if err := repo.db.Save(buyerOrder).Error; err != nil {
			http.Error(w, "Failed to save order", http.StatusInternalServerError)
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

	// Acknowledge the webhook request early
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Paystack Webhook received successfully")

	// Access the event field
	event := webhookPayload.Event
	data := webhookPayload.Data

	// Process the event
	switch event {
	case "charge.success":
		buyerOrder := &store.Order{}
		// Handle charge success logic here
		err := repo.db.Model(buyerOrder).Where("trt_ref = ? ", data.Reference).Error
		if err != nil {
			http.Error(w, "Failed to find order", http.StatusNotFound)
			return
		}
		buyerOrder.TransStatus = data.Status
		buyerOrder.Status = "pending"
		buyerOrder.PaymentMethod = data.Channel
		buyerOrder.PaymentGateway = "paystack"
		if err := repo.db.Save(buyerOrder).Error; err != nil {
			http.Error(w, "Failed to save order", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Received charge success event: %s for email: %s, Amount: %d\n", event, webhookPayload.Data.Customer.Email, webhookPayload.Data.Amount)

	case "charge.failed":
		// Handle charge failure logic here
		buyerOrder := &store.Order{}
		// Handle charge success logic here
		err := repo.db.Model(buyerOrder).Where("trt_ref = ? ", data.Reference).Error
		if err != nil {
			http.Error(w, "Failed to find order", http.StatusNotFound)
			return
		}
		buyerOrder.TransStatus = data.Status
		buyerOrder.Status = "pending"
		buyerOrder.PaymentMethod = data.Channel
		buyerOrder.PaymentGateway = "paystack"
		if err := repo.db.Save(buyerOrder).Error; err != nil {
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

func SquadWebhookHandler(w http.ResponseWriter, r *http.Request) {
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

	// Process the event
	switch event {
	case "payment.success":
		// Handle successful payment logic here
		fmt.Printf("Payment success: %s for email: %s, Amount: %d\n", webhookPayload.Data.Reference, webhookPayload.Data.Customer.Email, webhookPayload.Data.Amount)

	case "payment.failed":
		// Handle failed payment logic here
		fmt.Printf("Payment failed: %s for email: %s\n", webhookPayload.Data.Reference, webhookPayload.Data.Customer.Email)

	default:
		// Optionally handle other events or log them
		fmt.Println("Received unhandled event:", event)
	}

	// Log the received body for debugging or further processing
	fmt.Println("Received Webhook Body:", string(body))
}
