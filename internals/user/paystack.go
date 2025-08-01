package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// PaystackClient defines the interface for Paystack API operations
type PaystackClient interface {
	GetDVAAccount(email string) (*Account, error)
	CreateDVAAccount(details *DVADetails) (*Account, error)
	InitiateTransfer(ctx context.Context, req *TransferRequest) (*Transfer, error)
	CreateTransferRecipient(ctx context.Context, req *RecipientRequest) (*RecipientResponse, error)
	GetBanks(ctx context.Context) (*BanksResponse, error)
}

type TransferRequest struct {
	Amount    float64 `json:"amount"`
	Recipient string  `json:"recipient"`
	Reason    string  `json:"reason"`
}

type Transfer struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Reference     string  `json:"reference"`
		Amount        float64 `json:"amount"`
		Status        string  `json:"status"`
		TransferCode  string  `json:"transfer_code"`
		RecipientCode string  `json:"recipient_code"`
	} `json:"data"`
}

type RecipientRequest struct {
	Type          string `json:"type"`
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
}

type RecipientResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		RecipientCode string `json:"recipient_code"`
	} `json:"data"`
}

type PaystackBank struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Slug string `json:"slug"`
}

type BanksResponse struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    []PaystackBank `json:"data"`
}

// Add the implementation
func (p *paystackClient) InitiateTransfer(ctx context.Context, req *TransferRequest) (*Transfer, error) {
	url := "https://api.paystack.co/transfer"

	payload := map[string]interface{}{
		"amount":    req.Amount * 100, // Convert to kobo
		"recipient": req.Recipient,
		"reason":    req.Reason,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.secretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var transfer Transfer
	if err := json.NewDecoder(resp.Body).Decode(&transfer); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &transfer, nil
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
		"first_name":     details.StoreName,
		"last_name":      "",
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

// CreateTransferRecipient creates a new transfer recipient in Paystack
func (p *paystackClient) CreateTransferRecipient(ctx context.Context, req *RecipientRequest) (*RecipientResponse, error) {
	url := "https://api.paystack.co/transferrecipient"

	payload := map[string]interface{}{
		"type":           req.Type,
		"name":           req.Name,
		"account_number": req.AccountNumber,
		"bank_code":      req.BankCode,
		"currency":       "NGN",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.secretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var result RecipientResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if !result.Status {
		return nil, fmt.Errorf("paystack error: %s", result.Message)
	}

	return &result, nil
}

// GetBanks fetches all supported banks from Paystack
func (p *paystackClient) GetBanks(ctx context.Context) (*BanksResponse, error) {
	url := "https://api.paystack.co/bank"

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.secretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var result BanksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if !result.Status {
		return nil, fmt.Errorf("paystack error: %s", result.Message)
	}

	return &result, nil
}
