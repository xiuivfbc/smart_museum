package controller

import (
	"group_ten_server/dao"
	"group_ten_server/model"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key")

// UserClaims 用于JWT
type UserClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// isValidUserInput 校验用户名和密码为不包含空格的utf8
func isValidUserInput(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}
	if strings.ContainsAny(s, " ") {
		return false
	}
	// 只要是utf8即可
	return true
}

// Register 用户注册
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if !isValidUserInput(req.Username) || !isValidUserInput(req.Password) || strings.TrimSpace(req.Role) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名、密码和角色不能为空且不能包含空格"})
		return
	}
	// 检查用户名是否已存在
	exist, _ := dao.GetUserByUsername(req.Username)
	if exist != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}
	user := model.User{
		Username:  req.Username,
		Password:  req.Password, // 实际应加密
		Role:      req.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := dao.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if !isValidUserInput(req.Username) || !isValidUserInput(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码不能包含空格且不能为空"})
		return
	}
	user, err := dao.GetUserByUsername(req.Username)
	if err != nil || user == nil || user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}
	// 生成JWT
	claims := UserClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
