package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiuivfbc/smart_museum/service"
)

// Register 用户注册
func RegisterUser(c *gin.Context) {
	var req struct {
		Username         string `json:"username"`
		Password         string `json:"password"`
		Role             string `json:"role"`
		Phone            string `json:"phone"`
		Email            string `json:"email"`
		Identifier       string `json:"identifier"`
		VerificationCode string `json:"verification_code"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	service.RegisterUserService(c, req)
}

// Login 用户登录
func LoginUser(c *gin.Context) {
	var req struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	service.LoginUserService(c, req)
}

func UploadUser(c *gin.Context) {
	var req struct {
		Id          int    `json:"id"`
		Username    string `json:"username"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	service.UploadUserService(c, req)
}

func GetVerificationCode(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	service.GetVerificationCodeService(c, req)
}
