package middleware

import (
	"net/http"
	"strings"

	"go-blog-api/internal/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于 JWT 的认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1.从 HTTP 请求头中获取 Authorization 字段
		//标准的 Token 格式是放在 Header 里：Authorization：Bearer xxxxx.yyyyy.zzzzz
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请求头中缺失 Authorization 字段，请先登录"})
			c.Abort()
			return
		}
		//2. 按空格分割，提取出真正的 Token 字符串
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization 格式错误，应为 Bearer <Token>"})
			c.Abort()
			return
		}
		//parts[1] 就是我们要的 jwt Token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 无效或已过期，请重新登录"})
			c.Abort()
			return
		}
		//3. 验证通过！把提取出的 userID 存入 Gin 的上下文（Context）中
		//这样后续的业务接口（比如发文章）就能直接从Context里知道当前操作的是谁
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		//4.放行，继续执行后面的业务逻辑
		c.Next()

	}
}
