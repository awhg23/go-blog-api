package utils

import (
	"go-blog-api/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateAndParseToken 测试 JWT 的生成和解析逻辑是否闭环
func TestGenerateAndParseToken(t *testing.T) {
	// 1.准备工作：自定义一个配置，防止空指针报错，因为测试环境未执行 main.go 中的 InitConfig 函数
	config.App = &config.Config{
		JWT: config.JWTConfig{
			Secret: "test_super_secret_key",
		},
	}

	testUserID := int64(888)
	testUsername := "golang_tester"

	// 2.执行 Token 生成
	token, err := GenerateToken(testUserID, testUsername)

	// 断言（Assert）：预期 err 为 nil，且 token 不为空
	assert.NoError(t, err, "生成 Token 时不应返回错误")
	assert.NotEmpty(t, token, "生成的 Token 不应为空")

	// 3.执行解析 Token 测试
	claims, err := ParseToken(token)

	// 断言：解析不能报错，且解析出来的数据必须和存进去的一模一样！
	assert.NoError(t, err, "解析 Token 时不应返回错误")
	assert.Equal(t, testUserID, claims.UserID, "解析出的 UserID 不匹配")
	assert.Equal(t, testUsername, claims.Username, "解析出的 Username 不匹配")
}
