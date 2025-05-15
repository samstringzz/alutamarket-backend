package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// GetPaystackDVABalance gets the balance of a PayStack Dedicated Virtual Account
func GetPaystackDVABalance(accountNumber string) (float64, error) {
	url := "https://api.paystack.co/transaction"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return 0, err
	}

	// Use virtual_account_number instead of recipient_account
	q := req.URL.Query()
	q.Add("virtual_account_number", accountNumber)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	var response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []struct {
			Amount          float64 `json:"amount"`
			Status          string  `json:"status"`
			Currency        string  `json:"currency"`
			Channel         string  `json:"channel"`
			GatewayResponse string  `json:"gateway_response"`
			Metadata        struct {
				ReceiverAccountNumber string `json:"receiver_account_number"`
			} `json:"metadata"`
		} `json:"data"`
		Meta struct {
			TotalVolume float64 `json:"total_volume"`
		} `json:"meta"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("error decoding response: %v", err)
	}

	if !response.Status {
		return 0, fmt.Errorf("failed to get transactions: %s", response.Message)
	}

	// Use the total_volume from meta, which represents the total amount
	// Convert from kobo to naira
	balance := response.Meta.TotalVolume / 100

	return balance, nil
}
