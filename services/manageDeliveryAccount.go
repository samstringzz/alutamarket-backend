package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// PayDeliveryFund sends a request to the Paystack API to charge a customer and returns an error if something goes wrong.
func PayDeliveryFund(amount float32) error {
	url := "https://api.paystack.co/charge"

	// Create the request payload
	payload := map[string]interface{}{
		"email":  "alutamarket",
		"amount": amount,
		"bank": map[string]string{
			"code":           "058",
			"account_number": "0936231445",
		},
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	// Check if the response status code indicates success
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("Success:", string(body))
		return nil
	}

	// Return error with status code and response body
	return fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(body))
}
