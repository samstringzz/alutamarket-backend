package app

import (
	"log"
	"net/http"
	"os"

	"github.com/Chrisentech/aluta-market-api/internals/messages"
	"github.com/Chrisentech/aluta-market-api/internals/product"
	"github.com/Chrisentech/aluta-market-api/internals/user"
	"github.com/Chrisentech/aluta-market-api/router"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		log.Printf("Received message: %s", message)

		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func Start() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
		// Continue as env vars might be set in the system
	}

	// Initialize repositories and handlers with error checking
	userRepo := InitializePackage(UserPackage)
	if userRepo == nil {
		log.Fatal("Failed to initialize user repository - database connection failed")
	}

	userRepository, ok := userRepo.(user.Repository)
	if !ok {
		log.Fatal("Failed to convert user repository interface")
	}

	userSrvc := user.NewService(userRepository)
	userHandler := user.NewHandler(userSrvc)

	// Fix product repository type assertion
	productRepo := InitializePackage(ProductPackage)
	if productRepo == nil {
		log.Fatal("Failed to initialize product repository")
	}

	productRepository, ok := productRepo.(product.Repository)
	if !ok {
		log.Fatal("Failed to convert product repository interface")
	}

	productSrvc := product.NewService(productRepository)
	productHandler := product.NewHandler(productSrvc)

	// Initialize message handler
	messageRepo := messages.NewRepository()
	messageService := messages.NewService(messageRepo)
	messageHandler := messages.NewMessageHandler(messageService)

	// Verify message handler initialization
	if messageHandler == nil {
		log.Fatal("Failed to initialize message handler")
	}

	// Initialize router with all required services and handlers
	r := router.InitRouter(userHandler, productHandler, productSrvc, messageHandler)

	// Add a health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Get port from env or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
