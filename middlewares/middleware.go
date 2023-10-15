package middlewares

import (
	"net/http"
	"os"
	"strings"
	"fmt"
	"github.com/Chrisentech/aluta-market-api/errors"
	"github.com/dgrijalva/jwt-go"
)

// Middleware for stateless JWT authentication and authorization
// func AuthMiddleware( requiredRole string, next http.Handler) http.Handler {

// 	secretKey := os.Getenv("SECRET_KEY")
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         // Extract the JWT token from the request header
//         tokenString := r.Header.Get("Authorization")
//         if tokenString == "" {
// 			errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "You can't access this resource")
//             return
//         }

//         // Remove "Bearer " prefix from the token string
//         tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

//         // Parse the JWT token
//         token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//             return []byte(secretKey), nil
//         })
//         if err != nil || !token.Valid {
//             errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "You can't access this resource")
//             return
//         }

//         // Check if the user has the required role
//         claims, ok := token.Claims.(jwt.MapClaims)
//         if !ok {
//             errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "You can't access this resource")
//             return
//         }

//         userRole := claims["role"].(string)
//         if userRole != requiredRole {
//              errors.NewAppError(http.StatusForbidden, "FORBIDDEN", "You do not have the permission to be here")
//             return
//         }

//	        // If the token is valid and the user has the required role, proceed to the next handler
//	        next.ServeHTTP(w, r)
//	    })
//	}
func AuthMiddleware(requiredRole string, tokenString string) error {

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
	fmt.Println(token)
	if err != nil || !token.Valid {
		return errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "You can't access this resource")

	}

	// Check if the user has the required role
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "You can't access this resource")

	}

	// userRole := claims["role"].(string)
	// if userRole != requiredRole {
	// 	return errors.NewAppError(http.StatusForbidden, "FORBIDDEN", "You do not have the permission to be here")

	// }

	return nil
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

// func ProtectedRoute(role string, resolverFunc func(http.ResponseWriter, *http.Request)) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         // Check user's authorization role using AuthMiddleware
//         authorizedHandler := authMiddleware(os.Getenv("SECRET_KEY"), role, resolverFunc)
//         authorizedHandler.ServeHTTP(w, r)

//	        // Continue with the resolver or mutation logic
//	        resolverFunc(w, r)
//	    })
//	}
// func AuthMiddleware(requiredRole string) error {
// 	secretKey := os.Getenv("SECRET_KEY")
// 	tokenString := r.Header.Get("Authorization")
// 	if tokenString == "" {
// 		return errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "You can't access this resource")

// 	}

// 	// Remove "Bearer " prefix from the token string
// 	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

// 	// Parse the JWT token
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(secretKey), nil
// 	})
// 	if err != nil || !token.Valid {
// 		return errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "You can't access this resource")

// 	}

// 	// Check if the user has the required role
// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "You can't access this resource")

// 	}

// 	userRole := claims["role"].(string)
// 	if userRole != requiredRole {
// 		return errors.NewAppError(http.StatusForbidden, "FORBIDDEN", "You do not have the permission to be here")

// 	}
// 	return nil
// }
