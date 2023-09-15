package utils

import (
    "net/http"
    "strings"
    "github.com/dgrijalva/jwt-go"
	"github.com/Chrisentech/aluta-market-api/errors"

)

// Middleware for stateless JWT authentication and authorization
func AuthMiddleware(secretKey string, requiredRole string, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract the JWT token from the request header
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
			errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "You can't access this resource")
            return
        }

        // Remove "Bearer " prefix from the token string
        tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

        // Parse the JWT token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secretKey), nil
        })
        if err != nil || !token.Valid {
            errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "You can't access this resource")
            return
        }

        // Check if the user has the required role
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "You can't access this resource")
            return
        }

        userRole := claims["role"].(string)
        if userRole != requiredRole {
             errors.NewAppError(http.StatusForbidden, "FORBIDDEN", "You can't access this resource")
            return
        }

        // If the token is valid and the user has the required role, proceed to the next handler
        next.ServeHTTP(w, r)
    })
}

// Middleware for stateful basic authentication
func BasicAuthMiddleware(username, password string, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user, pass, ok := r.BasicAuth()
        if !ok || user != username || pass != password {
            errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "You can't access this resource")
            return
        }

        // If the credentials are valid, proceed to the next handler
        next.ServeHTTP(w, r)
    })
}