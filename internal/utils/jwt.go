package utils

import (
	"go-blog-api/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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

	// 使用配置里的 Secret
	secret := []byte(config.App.JWT.Secret)
	return token.SignedString(secret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否正确
		return []byte(config.App.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	// 验证Token是否有效，并提取CLaims的数据
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
