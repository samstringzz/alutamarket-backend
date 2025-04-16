package utils

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
)

type PaystackVerificationResponse struct {
    Status  bool `json:"status"`
    Data    struct {
        Status    string `json:"status"`
        Reference string `json:"reference"`
    } `json:"data"`
}

func VerifyPaystackTransaction(reference string) (*PaystackVerificationResponse, error) {
    url := fmt.Sprintf("https://api.paystack.co/transaction/verify/%s", reference)
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    secretKey := os.Getenv("PAYSTACK_SECRET_KEY")
    if os.Getenv("ENVIRONMENT") == "development" {
        secretKey = os.Getenv("PAYSTACK_SECRET_KEY_DEV")
    }

    req.Header.Set("Authorization", "Bearer "+secretKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var verificationResponse PaystackVerificationResponse
    if err := json.Unmarshal(body, &verificationResponse); err != nil {
        return nil, err
    }

    return &verificationResponse, nil
}