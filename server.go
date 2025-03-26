package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Chrisentech/aluta-market-api/graph"
	"github.com/Chrisentech/aluta-market-api/internals/messages"
	"github.com/Chrisentech/aluta-market-api/internals/product"
	"github.com/Chrisentech/aluta-market-api/internals/user"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	// "github.com/Chrisentech/aluta-market-api/app"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const defaultPort = "8080"

// WebSocket handler
type Message *messages.Message

// ExtractUserIDFromRequest extracts the user ID from the JWT in the Authorization header

func ExtractTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		ctx := r.Context()
		ctx = context.WithValue(ctx, "token", tokenString)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Update with the specific origin if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func InitServer() error {
	const uploadPath = "./uploads/"

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return fmt.Errorf("error creating upload directory: %v", err)
	}

	// Create a new CORS middleware with the desired options
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173", "https://www.thealutamarket.com", "https://thealutamarket.com", "https://alutamarket.vercel.app", "https://aluta-market-api-zns8.onrender.com/graphql"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize user components
	userRepo := user.NewRepository()
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	// Initialize product components
	productRepo := product.NewRepository()
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	// Initialize message components with proper error handling
	messageRepo := messages.NewRepository()
	if messageRepo == nil {
		return fmt.Errorf("failed to initialize message repository")
	}

	messageService := messages.NewService(messageRepo)
	if messageService == nil {
		return fmt.Errorf("failed to initialize message service")
	}

	messageHandler := messages.NewMessageHandler(messageService)
	if messageHandler == nil {
		return fmt.Errorf("failed to initialize message handler")
	}

	// Create resolver with all handlers using NewResolver
	resolver := graph.NewResolver(
		*userHandler,
		productService,
		productHandler,
		messageHandler,
	)

	// Debug log to verify resolver
	log.Printf("Resolver initialized successfully with MessageHandler: %+v", resolver.MessageHandler)

	// Configure the GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))

	// Add WebSocket transport with proper configuration
	// Configure the GraphQL server with secure WebSocket
	srv.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				return origin == "https://www.thealutamarket.com" ||
					origin == "https://thealutamarket.com" ||
					origin == "https://alutamarket.vercel.app" ||
					strings.HasPrefix(origin, "http://localhost")
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	// Set up routes
	router := gin.Default()

	// Apply CORS middleware
	router.Use(func(c *gin.Context) {
		corsMiddleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
	})

	// Add root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "running",
			"message": "Aluta Market API is running",
		})
	})

	// Add GraphQL playground endpoint
	router.GET("/graphql", gin.WrapH(playground.Handler("GraphQL Playground", "/graphql")))

	// GraphQL endpoint for queries/mutations
	router.POST("/graphql", gin.WrapH(srv))

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		if messageHandler == nil {
			log.Println("messageHandler is NIL when trying to handle WebSocket connection!")
		} else {
			log.Println("messageHandler is properly initialized, proceeding to handle WebSocket connection.")
			messageHandler.WebSocketHandler(c.Writer, c.Request)
		}
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// Bind to 0.0.0.0 to accept connections from any source
	return router.Run("0.0.0.0:" + port)
}
