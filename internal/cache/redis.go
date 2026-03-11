package cache

import (
	"context"
	"log"

	"go-blog-api/config"

	"github.com/redis/go-redis/v9"
)

// RDB 是全局的 Redis 客户端实例
var RDB *redis.Client

// InitReids 初始化 Redis 连接
func InitRedis() {
	cfg := config.App.Redis

	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}

	log.Println("Redis 连接成功!")
}
