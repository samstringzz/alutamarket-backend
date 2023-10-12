package utils

import (
	
	"github.com/google/uuid"
	"regexp"
	"strings"
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

// func MarshalJSON(data interface{}) []byte {
// 	jsonBytes, err := json.Marshal(data)
// 	if err != nil {
// 		return nil
// 	}

// 	return jsonBytes
// }

// func UnmarshalJSON(jsonBytes []byte, data interface{}) error {
// 	err := json.Unmarshal(jsonBytes, data)
// 	return err
// }


func GenerateUUID() string {
	uuid := uuid.New()
	return uuid.String()
}