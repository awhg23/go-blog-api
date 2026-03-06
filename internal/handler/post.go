package handler

import (
	"net/http"
	"strconv"

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

func GetPosts(c *gin.Context) {
	//1.获取分页参数（默认 page=1，size=10)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	//防御性校验：page和size不能小于1
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	} //防止单次请求过多导致崩溃

	offset := (page - 1) * size

	var posts []model.Post
	var total int64

	//2.获取文章总数（用于前段分页组件计算总页数)
	db.DB.Model(&model.Post{}).Count(&total)

	//3.分页查询并预加载作者信息
	if err := db.DB.Preload("User").
		Order("created_at desc").
		Limit(size).
		Offset(offset).
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, gin.H{
		"data": posts,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// UpdatePost 修改文章（需要 JWT鉴权）
func UpdatePost(c *gin.Context) {
	postID := c.Param("id")

	//1.获取当前登录的 userID （从 JWT 中间件中取得）
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var currentUserID int64
	switch v := userIDValue.(type) {
	case float64:
		currentUserID = int64(v)
	case int64:
		currentUserID = v
	case int:
		currentUserID = int64(v)
	}

	//2.绑定请求体参数
	var req struct {
		Title   string `json:"title" binding:"required,min=1,max=100"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不合法"})
		return
	}

	//3.校验越权：查询文章是否存在，并判断 user_id 是否匹配
	var post model.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	//【核心安全校验】：防止A修改B的文章
	if post.UserID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "越权操作：只能修改自己的文章"})
		return
	}

	//4.执行更新
	//使用 map 更新可以避免 GORM 忽略零值的问题，或者直接传入结构体
	if err := db.DB.Model(&post).Updates(map[string]interface{}{
		"title":   req.Title,
		"content": req.Content,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"data":    post,
	})
}

// DeletePost 删除文章 （需要 JWt 鉴权）
func DeletePost(c *gin.Context) {
	postID := c.Param("id")

	//1. 获取当前登录的 userID （从 JWT 中间件中取得）
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var currentUserID int64
	switch v := userIDVal.(type) {
	case float64:
		currentUserID = int64(v)
	case int64:
		currentUserID = v
	case int:
		currentUserID = int64(v)
	}
	//2.校验越权：查询文章并检查归属
	var post model.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	if post.UserID != currentUserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "越权操作，只能删除自己的文章"})
		return
	}

	//3.执行删除（这里是物理删除，想要软删除可以在 Model 中引入 gormDeletedAt）
	if err := db.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
