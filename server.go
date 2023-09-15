package main

import (
	"log"
	"net/http"
	"os"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Chrisentech/aluta-market-api/graph"
	"github.com/joho/godotenv"
    "github.com/Chrisentech/aluta-market-api/db"
    // "github.com/Chrisentech/aluta-market-api/app"
)

const defaultPort = "8080"

func main() {

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

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
	// Close the database connection when done
    // defer dbConnection.Close()
}

// go run github.com/99designs/gqlgen init
// go run github.com/99designs/gqlgen generate
