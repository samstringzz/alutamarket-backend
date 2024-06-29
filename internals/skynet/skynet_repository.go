package skynet

import (
	"bytes"
	"context"
	"encoding/json"
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

func (r *repository) BuyAirtime(ctx context.Context, input *Airtime) (*string, error) {
	// Assuming the GetUser function and user repository are defined elsewhere
	// customer, err := user.NewRepository().GetUser(ctx, input.UserID)
	// if err != nil {
	// 	return nil, errors.New("Customer not found")
	// }

	requestID := utils.GenerateRequestID()
	requestData := map[string]interface{}{
		"amount":     input.Amount,
		"request_id": requestID,
		"serviceID":  input.ServiceID,
		"phone":      input.Phone,
	}

	// Serialize the request data to JSON
	requestDataBytes, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	var requestURL, secretKey, apiKey string

	if os.Getenv("ENVIRONMENT") == "development" {
		requestURL = os.Getenv("VTU_DEV_URL")
		secretKey = os.Getenv("VTU_SANDBOX_SK")
		apiKey = os.Getenv("VTU_SANDBOX_API_KEY")
	} else {
		requestURL = os.Getenv("VTU_LIVE_URL")
		secretKey = os.Getenv("VTU_LIVE_SK")
		apiKey = os.Getenv("VTU_LIVE_API_KEY")
	}
	// fmt.Println(requestURL, secretKey, apiKey)

	// Create an HTTP client
	client := &http.Client{}

	// Create the HTTP request
	req, err := http.NewRequest("POST", requestURL+"/api/pay", bytes.NewBuffer(requestDataBytes))
	if err != nil {
		return nil, err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey)
	req.Header.Set("secret-key", secretKey)

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var transactionId, transactionStatus string

	// Check if the response status is not a success
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Parse the error response
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}

		// Print the entire error response to the console
		fmt.Printf("Error Response: %+v\n", errorResponse)

		// Return the error from the response if available
		if errMsg, ok := errorResponse["message"].(string); ok {
			return nil, errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", errMsg)
		}

		// If no specific error message is found, return a generic error
		return nil, errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", "unknown error occurred")
	}

	// Parse the response
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	// Extract transactionId and status from the successful response
	if content, ok := response["content"].(map[string]interface{}); ok {
		if transactions, ok := content["transactions"].(map[string]interface{}); ok {
			transactionId, _ = transactions["transactionId"].(string)
			transactionStatus, _ = transactions["status"].(string)
		}
	}

	// Print the entire response to the console
	fmt.Printf("Response: %+v\n", response)

	newService := &Skynet{
		UserID:        input.UserID,
		ID:            utils.GenerateUUID(),
		Status:        transactionStatus,
		RequestID:     requestID,
		Type:          "topup_airtime",
		Receiver:      input.Phone,
		TransactionID: transactionId,
	}

	err = r.db.Create(newService).Error
	if err != nil {
		return nil, err
	}

	successMsg := "Top up successful"
	return &successMsg, nil
}

func (r *repository) BuyData(ctx context.Context, input *Data) (*string, error) {

	requestID := utils.GenerateRequestID()
	requestData := map[string]interface{}{
		"amount":         input.Amount,
		"request_id":     requestID,
		"serviceID":      input.ServiceID,
		"phone":          input.Phone,
		"billersCode":    input.BillersCode,
		"variation_code": input.VariationCode,
	}

	// Serialize the request data to JSON
	requestDataBytes, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	var requestURL, secretKey, apiKey string

	if os.Getenv("ENVIRONMENT") == "development" {
		requestURL = os.Getenv("VTU_DEV_URL")
		secretKey = os.Getenv("VTU_SANDBOX_SK")
		apiKey = os.Getenv("VTU_SANDBOX_API_KEY")
	} else {
		requestURL = os.Getenv("VTU_LIVE_URL")
		secretKey = os.Getenv("VTU_LIVE_SK")
		apiKey = os.Getenv("VTU_LIVE_API_KEY")
	}
	// fmt.Println(requestURL, secretKey, apiKey)

	// Create an HTTP client
	client := &http.Client{}

	// Create the HTTP request
	req, err := http.NewRequest("POST", requestURL+"/api/pay", bytes.NewBuffer(requestDataBytes))
	if err != nil {
		return nil, err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey)
	req.Header.Set("secret-key", secretKey)

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var transactionId, transactionStatus string

	// Check if the response status is not a success
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Parse the error response
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}

		// Print the entire error response to the console
		fmt.Printf("Error Response: %+v\n", errorResponse)

		// Return the error from the response if available
		if errMsg, ok := errorResponse["message"].(string); ok {
			return nil, errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", errMsg)
		}

		// If no specific error message is found, return a generic error
		return nil, errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", "unknown error occurred")
	}

	// Parse the response
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	// Extract transactionId and status from the successful response
	if content, ok := response["content"].(map[string]interface{}); ok {
		if transactions, ok := content["transactions"].(map[string]interface{}); ok {
			transactionId, _ = transactions["transactionId"].(string)
			transactionStatus, _ = transactions["status"].(string)
		}
	}

	// Print the entire response to the console
	fmt.Printf("Response: %+v\n", response)

	newService := &Skynet{
		UserID:        input.UserID,
		ID:            utils.GenerateUUID(),
		Status:        transactionStatus,
		RequestID:     requestID,
		Type:          "data_bundle_purchase",
		Receiver:      input.Phone,
		TransactionID: transactionId,
	}

	err = r.db.Create(newService).Error
	if err != nil {
		return nil, err
	}

	successMsg := "Top up successful"
	return &successMsg, nil
}

func (r *repository) GetSubscriptionsBundles(ctx context.Context, serviceID string) (*DataBundle, error) {
	var requestURL, publicKey, apiKey string

	if os.Getenv("ENVIRONMENT") == "development" {
		requestURL = os.Getenv("VTU_DEV_URL")
		publicKey = os.Getenv("VTU_SANDBOX_PK")
		apiKey = os.Getenv("VTU_SANDBOX_API_KEY")
	} else {
		requestURL = os.Getenv("VTU_LIVE_URL")
		publicKey = os.Getenv("VTU_LIVE_PK")
		apiKey = os.Getenv("VTU_LIVE_API_KEY")
	}

	// Create an HTTP client
	client := &http.Client{}

	// Create the HTTP request
	req, err := http.NewRequest("GET", requestURL+"/api/service-variations?serviceID="+serviceID, nil)
	if err != nil {
		return nil, err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey)
	req.Header.Set("public-key", publicKey)

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status is not a success
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Parse the error response
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}

		// Print the entire error response to the console
		fmt.Printf("Error Response: %+v\n", errorResponse)

		// Return the error from the response if available
		if errMsg, ok := errorResponse["message"].(string); ok {
			return nil, errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", errMsg)
		}

		// If no specific error message is found, return a generic error
		return nil, errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", "unknown error occurred")
	}

	// Parse the response
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	// Extract the content from the response
	content, ok := response["content"].(map[string]interface{})
	if !ok {
		return nil, errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", "Invalid response format")
	}

	// Extract variations from the content
	variationsData, ok := content["varations"].([]interface{})
	if !ok {
		return nil, errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", "Invalid variations format")
	}

	// Convert variations to the expected type
	var variations []BundleVariation
	for _, v := range variationsData {
		variationMap, ok := v.(map[string]interface{})
		if !ok {
			return nil, errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", "Invalid variations format")
		}

		variation := BundleVariation{
			VariationCode:   variationMap["variation_code"].(string),
			Name:            variationMap["name"].(string),
			VariationAmount: variationMap["variation_amount"].(string),
			FixedPrice:      variationMap["fixedPrice"].(string),
		}
		variations = append(variations, variation)
	}

	// Create the response
	responseData := &DataBundle{
		ServiceName:    content["ServiceName"].(string),
		ServiceID:      content["serviceID"].(string),
		ConvinienceFee: content["convinience_fee"].(string),
		Variations:     variations,
	}

	return responseData, nil
}

func (r *repository) VerifySmartCard(ctx context.Context, serviceId, billersCode string) (*SmartcardVerificationResponse, error) {
	requestData := map[string]interface{}{
		"serviceID":   serviceId,
		"billersCode": billersCode,
	}

	// Serialize the request data to JSON
	requestDataBytes, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	var requestURL, secretKey, apiKey string

	if os.Getenv("ENVIRONMENT") == "development" {
		requestURL = os.Getenv("VTU_DEV_URL")
		secretKey = os.Getenv("VTU_SANDBOX_SK")
		apiKey = os.Getenv("VTU_SANDBOX_API_KEY")
	} else {
		requestURL = os.Getenv("VTU_LIVE_URL")
		secretKey = os.Getenv("VTU_LIVE_SK")
		apiKey = os.Getenv("VTU_LIVE_API_KEY")
	}

	// Create an HTTP client
	client := &http.Client{}

	// Create the HTTP request
	req, err := http.NewRequest("POST", requestURL+"/api/merchant-verify", bytes.NewBuffer(requestDataBytes))
	if err != nil {
		return nil, err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey)
	req.Header.Set("secret-key", secretKey)

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status is not a success
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Parse the error response
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}

		// Print the entire error response to the console
		fmt.Printf("Error Response: %+v\n", errorResponse)

		// Return the error from the response if available
		if errMsg, ok := errorResponse["message"].(string); ok {
			return nil, errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", errMsg)
		}

		// If no specific error message is found, return a generic error
		return nil, errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", "unknown error occurred")
	}

	// Parse the response
	var response SmartcardVerificationResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	// Return the parsed response
	return &response, nil
}
