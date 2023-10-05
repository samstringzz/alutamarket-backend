package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Chrisentech/aluta-market-api/db"
	"github.com/Chrisentech/aluta-market-api/graph"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	// "github.com/Chrisentech/aluta-market-api/app"
)

const defaultPort = "8080"

func main() {
	// Create a new CORS middleware with the desired options
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},         // Specify the allowed origins
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},        // Specify allowed HTTP methods
		AllowedHeaders:   []string{"Authorization", "Content-Type"}, // Specify allowed headers
		AllowCredentials: true,                                      // Allow credentials like cookies
	})
	
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	// Run auto migration
	db.Migrate()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	srv.AddTransport(&transport.Websocket{
		   Upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool {
                // Check against your desired domains here
                 return r.Host == "http://localhost:5173"
            },
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
        },
	})
	// Wrap the GraphQL handler with CORS middleware
	// srv.Use(middlewares.AuthMiddleware(os.Getenv("SECRET_KEY","")))
	wrappedHandler := c.Handler(srv)
	// log.Fatalf("Failed to token: %v", tokenString)

	http.Handle("/graphql", wrappedHandler)
	http.Handle("/", playground.Handler("Aluta Market playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	// Close the database connection when done
	// defer dbConnection.Close()
}

// go run github.com/99designs/gqlgen init
// go run github.com/99designs/gqlgen generate
