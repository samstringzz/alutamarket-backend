package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateSlug(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace spaces with hyphens
	name = strings.ReplaceAll(name, " ", "-")

	// Remove special characters using regex
	regex := regexp.MustCompile("[^a-z0-9-]")
	name = regex.ReplaceAllString(name, "")

	return name
}

func GenerateUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

func GenerateRequestID() string {
	// Set the location to Africa/Lagos
	location, err := time.LoadLocation("Africa/Lagos")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return ""
	}

	// Get the current time in the Africa/Lagos timezone
	now := time.Now().In(location)

	// Format the date and time as YYYYMMDDHHII
	dateTime := now.Format("200601021504")

	// Generate a random alphanumeric string
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 12
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}

	// Concatenate the dateTime and the random alphanumeric string
	requestID := dateTime + string(randomString)

	return requestID
}

// Define the structure for adding a new subscriber to OneSignal
type OneSignalSubscriber struct {
	AppID         string `json:"app_id"`
	Identifier    string `json:"identifier"`
	DeviceType    int    `json:"device_type"`               // Use 11 for Email
	EmailAuthHash string `json:"email_auth_hash,omitempty"` // Optional, for email authentication
}

func AddEmailSubscriber(email string) error {
	url := "https://onesignal.com/api/v1/players"
	subscriber := OneSignalSubscriber{
		AppID:      os.Getenv("ONE_SIGNAL_APP_ID"),
		Identifier: email,
		DeviceType: 11, // 11 represents Email
	}

	jsonData, err := json.Marshal(subscriber)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Println("Response body:", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	return nil
}

func GenerateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var sb strings.Builder
	sb.Grow(length)
	for i := 0; i < length; i++ {
		randomIndex := rng.Intn(len(charset))
		sb.WriteByte(charset[randomIndex])
	}
	return sb.String()
}
