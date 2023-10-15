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


func GenerateUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

