package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type OneSignalEmailRequest struct {
	AppID              string            `json:"app_id"`
	EmailSubject       string            `json:"email_subject"`
	TemplateID         string            `json:"template_id"`
	IncludeEmailTokens []string          `json:"include_email_tokens"`
	IncludeUnsubcribed bool              `json:"include_unsubscribed"`
	Contents           map[string]string `json:"custom_data"`
}

func SendEmail(templateID, subject string, to []string, contents map[string]string) error {
	url := "https://onesignal.com/api/v1/notifications"

	emailRequest := OneSignalEmailRequest{
		AppID:              os.Getenv("ONE_SIGNAL_APP_ID"),
		TemplateID:         templateID,
		EmailSubject:       subject,
		IncludeEmailTokens: to,
		IncludeUnsubcribed: true,
		Contents:           contents,
	}

	jsonData, err := json.Marshal(emailRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Basic "+os.Getenv("ONE_SIGNAL_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-200 response: %d, %s", resp.StatusCode, string(body))
	}

	// Log the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Println("Response body:", string(body))

	return nil
}
