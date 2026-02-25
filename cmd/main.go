package main

import (
	"log"

	"go-blog-api/internal/db"
	"go-blog-api/internal/router"
)

func main() {
	db.InitDB()
	r := router.SetupRouter()
	log.Println("服务器启动中，监听端口: 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
