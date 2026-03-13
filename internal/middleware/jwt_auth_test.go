package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-blog-api/config"
	"go-blog-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	// 设置 Gin 为测试模式，减少多余日志
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// 挂载要测试的中间件
	r.Use(JWTAuthMiddleware())

	// 写一个假的受保护的接口
	r.GET("/protected", func(c *gin.Context) {
		userID, _ := c.Get("userID")
		c.JSON(http.StatusOK, gin.H{"message": "success", "userID": userID})
	})
	return r
}

func TestJWTAuthMiddleware(t *testing.T) {
	// 初始化测试配置
	config.App = &config.Config{
		JWT: config.JWTConfig{Secret: "test_secret"},
	}

	router := setupTestRouter()

	t.Run("不带 Token 应该返回 401", func(t *testing.T) {
		// 1. 创建一个模拟的 HTTP 请求 （GET /protected）
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		// 2. 创建一个模拟的相应记录器（用于接收返回值）
		w := httptest.NewRecorder()

		// 3. 将请求发给 Gin 路由
		router.ServeHTTP(w, req)

		// 4. 断言：预期被中间件拦截，返回401
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("携带无效 Token 应该返回 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer fake_invalid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("携带有效的 Token 应该返回 200 并放行", func(t *testing.T) {
		// 生成一个合法的 Token
		validToken, _ := utils.GenerateToken(1024, "test_user")

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// 断言：预期中间件放行，成功到达接口，返回 200
		assert.Equal(t, http.StatusOK, w.Code)
		// 还可以进一步断言返回的 JSON 中包含正确的 userID
		assert.Contains(t, w.Body.String(), "1024")
	})
}
