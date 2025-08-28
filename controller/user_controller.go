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

var jwtKey []byte

func InitJwtKey(key string) {
	jwtKey = []byte(key)
} // UserClaims 用于JWT

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
		Username   string `json:"username"`
		Password   string `json:"password"`
		Role       string `json:"role"`
		Phone      string `json:"phone"`
		Email      string `json:"email"`
		Identifier string `json:"identifier"` // admin验证码
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if !isValidUserInput(req.Username) || !isValidUserInput(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名、密码不能为空且不能包含空格"})
		return
	}
	if strings.TrimSpace(req.Phone) == "" && strings.TrimSpace(req.Email) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "电话和邮箱至少填写一个"})
		return
	}
	if ok, _ := dao.GetActivationCode(req.Identifier); req.Role == "admin" && !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的管理员激活码"})
		return
	} else {
		dao.DeleteActivationCode(req.Identifier)
		dao.AddActivationCode()
	}
	// 检查电话唯一
	if req.Phone != "" {
		var user model.User
		db := dao.GetDB()
		if err := db.Where("phone = ?", req.Phone).First(&user).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "电话已被注册"})
			return
		}
	}
	// 检查邮箱唯一
	if req.Email != "" {
		var user model.User
		db := dao.GetDB()
		if err := db.Where("email = ?", req.Email).First(&user).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已被注册"})
			return
		}
	}
	user := model.User{
		Username:  req.Username,
		Password:  req.Password, // 实际应加密
		Role:      req.Role,
		Phone:     req.Phone,
		Email:     req.Email,
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
		Identifier string `json:"identifier"` // 电话或邮箱
		Password   string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if strings.TrimSpace(req.Identifier) == "" || !isValidUserInput(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "账号和密码不能为空且不能包含空格"})
		return
	}
	var user *model.User
	var err error
	if strings.Contains(req.Identifier, "@") {
		user, err = dao.GetUserByEmail(req.Identifier)
	} else {
		user, err = dao.GetUserByPhone(req.Identifier)
	}
	if err != nil || user == nil || user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "账号或密码错误"})
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
