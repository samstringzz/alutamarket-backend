package utils

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/samstringzz/alutamarket-backend/errors"
)

// GetUserIDFromContext extracts the user ID from the JWT token in the context
func GetUserIDFromContext(ctx context.Context) (uint32, error) {
	// Get token from context
	tokenString, ok := ctx.Value("token").(string)
	if !ok {
		return 0, errors.NewAppError(401, "UNAUTHORIZED", "No authorization token found")
	}

	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return 0, errors.NewAppError(401, "UNAUTHORIZED", "Invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.NewAppError(401, "UNAUTHORIZED", "Invalid token claims")
	}

	// Get user ID from claims
	userIDStr, ok := claims["id"].(string)
	if !ok {
		return 0, errors.NewAppError(401, "UNAUTHORIZED", "Invalid user ID in token")
	}

	// Convert string ID to uint32
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0, errors.NewAppError(401, "UNAUTHORIZED", "Invalid user ID format")
	}

	return uint32(userID), nil
}
