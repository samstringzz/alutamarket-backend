package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Chrisentech/aluta-market-api/db"
	"github.com/Chrisentech/aluta-market-api/graph"
	"github.com/Chrisentech/aluta-market-api/internals/messages"
	"github.com/Chrisentech/aluta-market-api/services"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	// "github.com/Chrisentech/aluta-market-api/app"
)

const defaultPort = "8080"

// WebSocket handler
type Message *messages.Message

func ExtractTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		ctx := r.Context()
		ctx = context.WithValue(ctx, "token", tokenString)
		// Use the updated context when calling the next handler.

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	// Start broadcasting messages to connected clients
	go messages.BroadcastMessages()

	// Create a new CORS middleware with the desired options
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173", "https://www.thealutamarket.com", "https://thealutamarket.com", "https://thealutamarket.netlify.app"}, // Specify the allowed origins
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},                                                                                                                               // Specify allowed HTTP methods
		AllowedHeaders:   []string{"Authorization", "Content-Type"},                                                                                                                        // Specify allowed headers
		AllowCredentials: true,                                                                                                                                                             // Allow credentials like cookies
	})

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	repo := services.NewRepository()
	http.HandleFunc("/webhook/fw", repo.FWWebhookHandler)
	http.HandleFunc("/webhook/ps", repo.PaystackWebhookHandler)
	// http.HandleFunc("/webhook/squad", repo.SquadWebhookHandler)
	// Run auto migration
	db.Migrate()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return r.Host == os.Getenv("BASE_URL") || r.Host == "http://localhost:5173"
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	// Extract the token from the request context

	// Wrap the GraphQL handler with CORS middleware
	wrappedHandler := c.Handler(srv)
	http.Handle("/graphql", ExtractTokenMiddleware(wrappedHandler))
	http.Handle("/", playground.Handler("Aluta Market playground", "/query"))
	http.Handle("/query", ExtractTokenMiddleware(wrappedHandler))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Close the database connection when done
	// defer dbConnection.Close()
}

// go get github.com/99designs/gqlgen
// go run github.com/99designs/gqlgen generate
