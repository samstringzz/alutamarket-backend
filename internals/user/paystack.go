package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// PaystackClient defines the interface for Paystack API operations
type PaystackClient interface {
	GetDVAAccount(email string) (*Account, error)
	CreateDVAAccount(details *DVADetails) (*Account, error)
}

// paystackClient implements PaystackClient interface
type paystackClient struct {
	secretKey string
}

// NewPaystackClient creates a new Paystack client
func NewPaystackClient(secretKey string) PaystackClient {
	return &paystackClient{
		secretKey: secretKey,
	}
}

// Implement the interface methods
func (p *paystackClient) GetDVAAccount(email string) (*Account, error) {
	url := "https://api.paystack.co/dedicated_account"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	req.Header.Set("Authorization", "Bearer "+p.secretKey)
	req.Header.Set("Content-Type", "application/json")

	// Add query parameters
	q := req.URL.Query()
	q.Add("email", email)
	req.URL.RawQuery = q.Encode()

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse response
	var result struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Account Account `json:"account"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if !result.Status {
		return nil, fmt.Errorf("paystack error: %s", result.Message)
	}

	return &result.Data.Account, nil
}

func (p *paystackClient) CreateDVAAccount(details *DVADetails) (*Account, error) {
	url := "https://api.paystack.co/dedicated_account"

	payload := map[string]interface{}{
		"email":          details.StoreEmail,
		"first_name":     details.User.Fullname,
		"last_name":      details.StoreName,
		"phone":          details.User.Phone,
		"preferred_bank": "wema-bank",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	req.Header.Set("Authorization", "Bearer "+p.secretKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse response
	var result struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Account Account `json:"account"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if !result.Status {
		return nil, fmt.Errorf("paystack error: %s", result.Message)
	}

	return &result.Data.Account, nil
}
