package handler

import (
	"net/http"

	"go-blog-api/internal/db"
	"go-blog-api/internal/model"

	"github.com/gin-gonic/gin"
)

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=100"`
	Content string `json:"content" binding:"required"`
}

func CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误（标题和内容不为空）"})
		return
	}
	//核心逻辑：从上下文中获取刚刚中间件塞进去的userID
	userIDvalue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}
	//断言为int64
	userID := userIDvalue.(int64)

	//构建要插入数据库的文章模型
	post := model.Post{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}
	//存入数据库
	if err := db.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文章发布成功",
		"data": gin.H{
			"post_id": post.ID,
			"title":   post.Title,
		},
	})
}
