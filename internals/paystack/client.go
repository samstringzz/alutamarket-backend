package paystack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client defines the interface for Paystack API operations
type Client interface {
	GetDVAAccount(email string) (*Account, error)
	CreateDVAAccount(details *DVADetails) (*Account, error)
}

// Account represents a Paystack DVA account
type Account struct {
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	Email         string `json:"email"`
	Bank          struct {
		Name string `json:"name"`
	} `json:"bank"`
}

// DVADetails represents the details needed to create a DVA account
type DVADetails struct {
	StoreEmail string
	User       struct {
		Fullname string
		Phone    string
	}
	StoreName string
}

// paystackClient implements Client interface
type paystackClient struct {
	secretKey string
}

// NewClient creates a new Paystack client
func NewClient(secretKey string) Client {
	return &paystackClient{
		secretKey: secretKey,
	}
}

// GetDVAAccount retrieves a DVA account from Paystack
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
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []Account `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if !result.Status {
		return nil, fmt.Errorf("paystack error: %s", result.Message)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no DVA account found for email: %s", email)
	}

	// Find the account that matches the email
	for _, account := range result.Data {
		if account.Email == email {
			return &account, nil
		}
	}

	return nil, fmt.Errorf("no matching DVA account found for email: %s", email)
}

// CreateDVAAccount creates a new DVA account in Paystack
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
