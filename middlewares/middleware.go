package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin" // Add this import
	"github.com/golang-jwt/jwt/v4"
	"github.com/samstringzz/alutamarket-backend/errors"
)

func AuthMiddleware(requiredRole string, tokenString string) *errors.AppError {
	secretKey := os.Getenv("SECRET_KEY")

	if tokenString == "" {
		return errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "No authorization token passed!")
	}

	// Remove "Bearer " prefix from the token string
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return errors.NewAppError(http.StatusUnauthorized, "BAD REQUEST", "Invalid or expired token")
	}

	// Check if the user has the required role
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Invalid claims in the token")
	}

	userRole, ok := claims["usertype"].(string)
	if !ok {
		return errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user role claim in the token")
	}

	//configure role for both seller and buyer endpoints
	if userRole != requiredRole && userRole != "admin" && requiredRole != "entry" {
		return errors.NewAppError(http.StatusForbidden, "FORBIDDEN", "You do not have the permission to access this resource")
	}

	return nil
}

// Middleware for stateful basic authentication
func BasicAuthMiddleware(username, password string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Invalid basic authentication credentials")
			return
		}

		// If the credentials are valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// Update the import in router.go to use this ErrorHandler
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			var statusCode int
			var message string

			if appErr, ok := err.Err.(*errors.AppError); ok {
				statusCode = appErr.StatusCode
				message = appErr.Message
			} else {
				statusCode = http.StatusInternalServerError
				message = "Internal Server Error"
			}

			c.JSON(statusCode, gin.H{
				"error": message,
			})
			return
		}
	}
}
