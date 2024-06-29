package utils

import (
	"fmt"
	"math/rand"
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
