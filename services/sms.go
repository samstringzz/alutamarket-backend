package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type TermiiSMSRequest struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Sms     string `json:"sms"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	ApiKey  string `json:"api_key"`
}

type TermiiSMSResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

func SendSMS(to, from, message string) (*TermiiSMSResponse, error) {
	url := "https://api.ng.termii.com/api/sms/send"

	apiKey := os.Getenv("TERMII_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TERMII_API_KEY is not set")
	}

	// Prepare the SMS request payload
	smsRequest := TermiiSMSRequest{
		To:      to,
		From:    from,
		Sms:     message,
		Type:    "plain", // Use "plain" for plain text messages
		Channel: "dnd",   // Use "generic" for regular SMS; "dnd" for Do Not Disturb
		ApiKey:  apiKey,
	}

	// Log the request payload
	fmt.Printf("Sending SMS request: %+v\n", smsRequest)

	jsonData, err := json.Marshal(smsRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Make the HTTP POST request to Termii API
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Log the response body
	fmt.Printf("Received response: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d, %s", resp.StatusCode, string(body))
	}

	// Parse the response from Termii API
	var smsResponse TermiiSMSResponse
	if err := json.Unmarshal(body, &smsResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return &smsResponse, nil
}
