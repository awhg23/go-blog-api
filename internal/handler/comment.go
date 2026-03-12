package handler

import (
	"net/http"
	"strconv"

	"go-blog-api/internal/db"
	"go-blog-api/internal/model"
	"go-blog-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,max=500"`
}

// ==================== 1. 发表评论 (需要 JWT 鉴权) ====================
// CreateComment 发表评论
// @Summary 发表评论
// @Description 发表评论（需要登录）
// @Tags 评论模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateCommentRequest true "评论内容"
// @Param id path int true "文章ID"
// @Success 200 {object} map[string]interface{} "评论成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未登录或Token无效"
// @Failure 404 {object} map[string]interface{} "文章不存在"
// @Failure 500 {object} map[string]interface{} "发表评论失败"
// @Router /posts/:id/comments [post]
func CreateComment(c *gin.Context) {
	// 1. 获取 URL 路径中的文章 ID，例如 /api/posts/1/comments
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	//2.检查文章是否存在（防御性编程）
	var post model.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	//3.从上下文获取当前登录的用户 ID （复用之前的 switch 断言逻辑）
	currentUserID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或Token无效"})
		return
	}

	//4.解析前端传来的评论内容
	var req struct {
		Content string `json:"content" binding:"required,max=500"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论内容不能为空且不能超过500字"})
		return
	}

	//5.组装并入库
	comment := model.Comment{
		PostID:  postID,
		UserID:  currentUserID,
		Content: req.Content,
	}
	if err := db.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发表评论失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "评论成功",
		"data":    comment,
	})
}

// ==================== 2. 获取某篇文章的评论列表 (公开接口) ====================
// GetPostComments 获取某篇文章的评论列表
// @Summary 获取某篇文章的评论列表（分页）
// @Description 获取某篇文章的评论列表（公开接口）
// @Tags 评论模块
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param page query int false "页码，默认1"
// @Param size query int false "每页数量，默认10"
// @Success 200 {object} map[string]interface{} "返回评论列表和分页信息"
// @Failure 404 {object} map[string]interface{} "文章不存在"
// @Failure 500 {object} map[string]interface{} "获取评论列表失败"
// @Router /posts/:id/comments [get]
func GetPostComments(c *gin.Context) {
	postID := c.Param("id")

	//检查文章是否存在（防御性编程）
	var post model.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 获取分页参数（默认 page=1，size=10）
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	//防御性校验： page和size不能小于1
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	offset := (page - 1) * size

	var total int64
	var comments []model.Comment

	//获取评论数量（用于前端分页组件计算总页数）
	db.DB.Model(&model.Comment{}).
		Where("post_id = ?", postID).
		Count(&total)

	// 按照创建时间倒序排列（最新的评论在最上面），并且预加载评论的作者信息
	if err := db.DB.Preload("User").
		Where("post_id = ?", postID).
		Order("created_at desc").
		Limit(size).
		Offset(offset).
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论列表失败"})
		return
	}

	//返回结果
	c.JSON(http.StatusOK, gin.H{
		"data": comments,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// ==================== 3. 删除评论 (需要 JWT 鉴权) ====================
// DeleteComment 删除评论
// @Summary 删除评论
// @Description 删除评论（需要登录，只能删除自己的评论）
// @Tags 评论模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "评论ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 401 {object} map[string]interface{} "未登录或Token无效"
// @Failure 403 {object} map[string]interface{} "越权操作：只能删除自己的评论"
// @Failure 404 {object} map[string]interface{} "评论不存在"
// @Failure 500 {object} map[string]interface{} "删除失败"
// @Router /comments/:id [delete]
func DeleteComment(c *gin.Context) {
	commentID := c.Param("id")

	//1.获取当前登录用户ID
	currentUserID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或Token无效"})
		return
	}

	//2.查找评论
	var comment model.Comment
	if err := db.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	//3.越权校验：只能删除自己的评论
	if comment.UserID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "越权操作：只能删除自己的评论"})
		return
	}

	//4.执行删除
	if err := db.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评论已删除"})
}
