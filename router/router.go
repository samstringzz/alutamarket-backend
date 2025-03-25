package router

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Chrisentech/aluta-market-api/graph"
	"github.com/Chrisentech/aluta-market-api/internals/messages"
	"github.com/Chrisentech/aluta-market-api/internals/product"
	"github.com/Chrisentech/aluta-market-api/internals/user"
	"github.com/Chrisentech/aluta-market-api/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

// Update the InitRouter function to accept both handlers
// Update function signature to include productService
func InitRouter(userHandler *user.Handler, productHandler *product.Handler, productService product.Service, messageHandler *messages.Handler) *gin.Engine {
	r := gin.Default()

	// Add error handler middleware
	r.Use(middlewares.ErrorHandler())

	// CORS middleware configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:5173",
		"http://localhost:3000",
		"https://aluta-market-api-zns8.onrender.com",
		"https://*.onrender.com",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Requested-With",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Methods",
	}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}
	r.Use(cors.New(config))

	// Update GraphQL server configuration to include both handlers and service
	// Ensure message handler is passed to GraphQL server configuration
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			UserHandler:    *userHandler,
			ProductHandler: productHandler,
			ProductService: productService,
			MessageHandler: messageHandler,
		},
	}))
	r.POST("/graphql", gin.WrapH(srv))
	r.GET("/graphql", gin.WrapH(playground.Handler("GraphQL playground", "/graphql")))

	// Basic routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Aluta Market API",
		})
	})

	// Add product routes
	// Remove the REST endpoints as we're using GraphQL
	// productRoutes := r.Group("/api/products")
	// {
	//     productRoutes.POST("/category", productHandler.CreateCategory)
	//     productRoutes.GET("/categories", productHandler.GetCategories)
	// }

	// Add WebSocket route
	r.GET("/ws", func(c *gin.Context) {
		if messageHandler == nil {
			log.Println("Message handler is nil!")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		messageHandler.WebSocketHandler(c.Writer, c.Request)
	})

	return r
}

func Start(addr string) error {
	return r.Run(addr)
}
