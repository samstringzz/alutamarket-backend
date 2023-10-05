package router

import (
	"fmt"
	"net/http"
	"os"
	"github.com/Chrisentech/aluta-market-api/errors"
	"github.com/golang-jwt/jwt/v4"
    "context"
)

func isValidSession(tokenString string, secretKey string) bool {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Validate the signing method and return the secret key for verification
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("invalid signing method")
        }
        return []byte(secretKey), nil
    })
    if err != nil || !token.Valid {
        return false
    }
    return true
}

// Define a custom middleware function that checks the authentication.
func AuthMiddleware(ctx context.Context) (context.Context, error) {
    // Extract the token from the context or request
    token, ok := ctx.Value("token").(string)
    if !ok {
        return ctx, errors.NewAppError(http.StatusForbidden, "UNAUTHORIZED", "Authorization token not provided")
    }

    // Validate the token using your token validation logic
    if !isValidSession(token, os.Getenv("SECRET_KEY")) {
        return ctx, errors.NewAppError(http.StatusForbidden, "UNAUTHORIZED", "Invalid or expired token")
    }

    // If the token is valid, proceed with the resolver function
    return ctx, nil
}
