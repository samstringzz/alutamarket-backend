package router

import (
	"fmt"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Retrieve the cookie
        cookie, err := c.Request.Cookie("cookie-session")
        if err != nil {
            // Cookie not found or error occurred
            // Handle the error or redirect as needed
            c.Redirect(http.StatusSeeOther, "/login")
            return
        }

        // Extract the value from the cookie
        cookieValue := cookie.Value

        // Perform authorization checks (e.g., validate access token or user ID)
        if !isValidSession(cookieValue,os.Getenv("SECRET_KEY")) {
            // Access denied
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        // Access granted, continue to the next middleware or handler
        c.Next()
    }
}
