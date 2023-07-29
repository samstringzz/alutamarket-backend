package router

import (
	"github.com/Chrisentech/aluta-market-api/internals/user"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler) {
	r = gin.Default()

	r.POST("/register", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.GET("/logout", userHandler.Logout)

	r.POST("/product", AuthMiddleware(), userHandler.CreateUser)
}

func Start(addr string) error {
	return r.Run(addr)
}
