package router

import (
	// "net/http"
	"github.com/Chrisentech/aluta-market-api/internals/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler) {
	r = gin.Default()
	// CORS middleware configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))

	// r.GET("/", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "Hello Server")
	// })
	// r.POST("/register", userHandler.CreateUser)
	// r.POST("/login", userHandler.Login)
	// r.GET("/logout", userHandler.Logout)

	// r.POST("/product", AuthMiddleware(), userHandler.CreateUser)
}

func Start(addr string) error {
	return r.Run(addr)
}
