package router

import (
	"go-blog-api/internal/handler"
	"go-blog-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)

		authApi := api.Group("")
		authApi.Use(middleware.JWTAuthMiddleware())
		{
			authApi.POST("/posts", handler.CreatePost)
		}
	}

	return r
}
