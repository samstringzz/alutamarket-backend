package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRandKey() (string,error) {
	// Generate a random byte slice with sufficient length
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println("Error generating secret key:", err)
		return "",err
	}

	// Encode the byte slice as a base64 string
	secretKey := base64.URLEncoding.EncodeToString(key)
	fmt.Println("Generated secret key:", secretKey)
	return secretKey,nil
}






