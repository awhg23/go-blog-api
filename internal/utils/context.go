package utils

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCurrentUserID(c *gin.Context) (int64, error) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或Token无效"})
		return 0, errors.New("未登录或Token无效")
	}

	switch v := userIDValue.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "上下文中无效的userID类型"})
		return 0, errors.New("上下文中无效的userID类型")
	}
}
