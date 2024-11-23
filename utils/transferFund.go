package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// PayFund handles the process of creating a Paystack transfer recipient and initiating a transfer.
func PayFund(amount float32, accountNumber, bankCode string) error {
	amountInKobo := int(amount * 100)
	if amountInKobo > 0 {
		// 1. Create the transfer recipient
		recipientCode, err := createTransferRecipient(accountNumber, bankCode)
		if err != nil {
			return fmt.Errorf("error creating transfer recipient: %w", err)
		}

		// 2. Initiate the transfer using the recipient code
		err = initiateTransfer(amountInKobo, recipientCode)
		if err != nil {
			return fmt.Errorf("error initiating transfer: %w", err)
		}

		return nil
	} else {
		return fmt.Errorf("invalid amount: %f", amount)
	}
}

// createTransferRecipient creates a new transfer recipient in Paystack
func createTransferRecipient(accountNumber, bankCode string) (string, error) {
	url := "https://api.paystack.co/transferrecipient"

	// Payload for creating a transfer recipient
	payload := map[string]interface{}{
		"type":           "nuban",
		"name":           "Aluta Logistics", // You can adjust this or pass the name as a parameter
		"account_number": accountNumber,
		"bank_code":      bankCode,
		"currency":       "NGN", // Assuming NGN. Adjust based on your needs.
		"email":          "folajimiopeyemisax3@gmail.com",
	}

	// Marshal payload into JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshalling payload: %w", err)
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Set request headers
	req.Header.Set("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read and check response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Check for success
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Parse the response to extract the recipient code
		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return "", fmt.Errorf("error unmarshalling response: %w", err)
		}

		data := response["data"].(map[string]interface{})
		recipientCode := data["recipient_code"].(string)
		return recipientCode, nil
	}

	// Return error with status code and response body if failed
	return "", fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(body))
}

// initiateTransfer sends money to the recipient created using Paystack transfer API
func initiateTransfer(amountInKobo int, recipientCode string) error {
	url := "https://api.paystack.co/transfer"

	// Payload for initiating a transfer
	payload := map[string]interface{}{
		"source":    "balance",
		"amount":    amountInKobo,
		"recipient": recipientCode,
		"reason":    "Wthdrawal fund", // Customize this message
	}

	// Marshal payload into JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload: %w", err)
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set request headers
	req.Header.Set("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read and check response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	// Check for success
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("Transfer Successful:", string(body))
		return nil
	}

	// Return error with status code and response body if failed
	return fmt.Errorf("transfer failed with status code %d: %s", resp.StatusCode, string(body))
}
