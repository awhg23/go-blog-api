package db

import (
	"fmt"
	"go-blog-api/internal/model"
	"log"

	"go-blog-api/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// 从全局配置中获取数据库信息
	dbCfg := config.App.Database

	// DSN (Data Source Name): 账号:密码@tcp(地址:端口)/数据库名?参数
	// 动态拼接DNS
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	// 自动迁移（根据模型创建表）
	DB.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	log.Println("数据库连接成功")
}
