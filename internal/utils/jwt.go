package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret 是用于加密和解密 Token 的密钥（类似保险柜的钥匙）
// ⚠️ 注意：在真实的企业项目中，这个密钥千万不能写死在代码里，应该从 config.yaml 或环境变量读取！
var jwtSecret = []byte("my_super_secret_key_123456")

// Claims 定义我们要在 Token 里携带的信息（载荷）
type Claims struct {
	UserID               int64  `json:"user_id"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        //JWT官方内置的字段 ？？
}

func GenerateToken(userID int64, username string) (string, error) {
	//1.设置Token过期时间
	expirationTime := time.Now().Add(24 * time.Hour)
	//组装要装进Token的数据
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), //过期时间
			Issuer:    "go-blog-api",                      //签发人
		},
	}

	//使用 HS256 算法生成Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
