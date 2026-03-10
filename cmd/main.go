package main

import (
	"log"

	"go-blog-api/config"
	_ "go-blog-api/docs"
	"go-blog-api/internal/db"
	"go-blog-api/internal/router"
)

// @title Go 博客后端 API
// @version 1.0
// @description 这是一个基于 Go + Gin + GORM 开发的博客项目后端 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email awhg23@outlook.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// [新增] 1.初始化配置
	config.InitConfig()

	// 2.初始化数据库连接（此时 db.go 中才能安全读取到 config.App.Database)
	db.InitDB()

	// 3.初始化路由
	r := router.SetupRouter()

	// 4.从配置中读取端口并启动服务器
	port := ":" + config.App.Server.Port
	log.Printf("服务器启动中，监听端口 %s...", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
