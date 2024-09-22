package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type PaystackResponse struct {
	Status bool `json:"status"`
	Data   struct {
		AccountName string `json:"account_name"`
	} `json:"data"`
	Message string `json:"message"`
}

func verifyAccountNumber(bankName, accountNumber string) (string, error) {
	KEY := os.Getenv("PAYSTACK_SECRET_KEY") // replace with your actual Paystack secret key

	if bankName != "" && accountNumber != "" {
		client := &http.Client{}
		url := fmt.Sprintf("https://api.paystack.co/bank/resolve?account_number=%s&bank_code=%s", accountNumber, bankName)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return "", fmt.Errorf("error creating request: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+KEY)

		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("error sending request: %v", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error reading response body: %v", err)
		}

		var result PaystackResponse
		if err := json.Unmarshal(body, &result); err != nil {
			return "", fmt.Errorf("error unmarshalling response: %v", err)
		}

		if result.Status {
			return result.Data.AccountName, nil
		} else {
			return "", fmt.Errorf("account verification failed: %s", result.Message)
		}
	}

	return "", fmt.Errorf("bank name or account number is missing")
}

// Handler function for HTTP requests
func VerifyAccountNumberHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	bankName := r.URL.Query().Get("bankCode")
	accountNumber := r.URL.Query().Get("accountNumber")

	// Call the VerifyAccountNumber function
	accountName, err := verifyAccountNumber(bankName, accountNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"accountName": accountName,
	}
	json.NewEncoder(w).Encode(response)
}
