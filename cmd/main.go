package main

import (
	"log"

	"go-blog-api/config"
	"go-blog-api/internal/db"
	"go-blog-api/internal/router"
)

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
