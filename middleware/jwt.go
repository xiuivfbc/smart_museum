package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

type UserClaims struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func InitJwtKey(key string) {
	jwtKey = []byte(key)
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未携带token"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token无效"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*UserClaims); ok {
			// 验证是否为 AccessToken
			if claims.TokenType != "access" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token类型错误"})
				c.Abort()
				return
			}

			// 将用户信息存储到上下文中
			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token解析失败"})
			c.Abort()
			return
		}

		c.Next()
	}
}
