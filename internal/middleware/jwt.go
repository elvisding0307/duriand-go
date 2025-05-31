package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization header 获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		// 解析和验证 token
		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			// 确保使用的是正确的签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte("your-secret-key"), nil // 注意：实际使用时应该从配置文件读取密钥
		})

		if err != nil || !token.Valid {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		// 将 token 中的信息存储到上下文中
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			uid := uint64(claims["uid"].(float64))
			c.Set("uid", uid)
		}

		c.Next()
	}
}
