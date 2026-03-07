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
		// ========= 公开接口（所有人可访问） =========
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
		api.GET("/posts", handler.GetPosts)

		api.GET("/posts/:id/comments", handler.GetPostComments)
		// ========= 私有接口（需要 JWT Token 鉴权） =========
		authApi := api.Group("")
		authApi.Use(middleware.JWTAuthMiddleware())
		{
			authApi.POST("/posts", handler.CreatePost)
			authApi.PUT("/posts/:id", handler.UpdatePost)
			authApi.DELETE("/posts/:id", handler.DeletePost)
			authApi.POST("/posts/:id/comments", handler.CreateComment)
			authApi.DELETE("/comments/:id", handler.DeleteComment)
		}
	}

	return r
}
