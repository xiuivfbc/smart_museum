package service

import (
	"errors"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/xiuivfbc/smart_museum/config"
	"github.com/xiuivfbc/smart_museum/dao"
	"github.com/xiuivfbc/smart_museum/model"
	"gopkg.in/gomail.v2"
)

var jwtKey []byte
var RedisClient *redis.Client

type UserClaims struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"` // "access" 或 "refresh"
	jwt.RegisteredClaims
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func InitJwtKey(key string) {
	jwtKey = []byte(key)
}

// 生成 AccessToken (短期有效，15分钟)
func generateAccessToken(userID int, username, role string) (string, error) {
	claims := UserClaims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// 生成 RefreshToken (长期有效，7天)
func generateRefreshToken(userID int, username, role string) (string, error) {
	claims := UserClaims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// 验证并解析 Token
func validateToken(tokenString string, expectedType string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		if claims.TokenType != expectedType {
			return nil, errors.New("token type mismatch")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func RegisterUserService(c *gin.Context, req interface{}) {
	reqStruct, ok := req.(struct {
		Username         string `json:"username"`
		Password         string `json:"password"`
		Role             string `json:"role"`
		Phone            string `json:"phone"`
		Email            string `json:"email"`
		Identifier       string `json:"identifier"`
		VerificationCode string `json:"verification_code"`
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数类型错误"})
		return
	}
	if !isValidUserInput(reqStruct.Password) {
		c.JSON(http.StatusOK, gin.H{"error": "密码不能为空且不能包含空格"})
		return
	}
	if strings.TrimSpace(reqStruct.Phone) == "" && strings.TrimSpace(reqStruct.Email) == "" {
		c.JSON(http.StatusOK, gin.H{"error": "电话和邮箱至少填写一个"})
		return
	}
	if reqStruct.Email != "" && !IsValidEmail(reqStruct.Email) {
		c.JSON(http.StatusOK, gin.H{"error": "邮箱格式不正确"})
		return
	}
	if reqStruct.Phone != "" {
		var user model.User
		db := dao.GetDB()
		if err := db.Where("phone = ?", reqStruct.Phone).First(&user).Error; err == nil {
			c.JSON(http.StatusOK, gin.H{"error": "电话已被注册"})
			return
		}
	}
	if reqStruct.Email != "" {
		var user model.User
		db := dao.GetDB()
		if err := db.Where("email = ?", reqStruct.Email).First(&user).Error; err == nil {
			c.JSON(http.StatusOK, gin.H{"error": "邮箱已被注册"})
			return
		}
	}
	// 管理员注册需要激活码验证
	if reqStruct.Role == "admin" {
		if ok, _ := dao.GetActivationCode(reqStruct.Identifier); !ok {
			c.JSON(http.StatusOK, gin.H{"error": "无效的管理员激活码"})
			return
		}
	}
	// if err := checkVerificationCode(c, reqStruct.Email, reqStruct.VerificationCode); err != nil {
	// 	c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	// 	return
	// }
	user := model.User{
		Username:  reqStruct.Username,
		Password:  reqStruct.Password,
		Role:      reqStruct.Role,
		Phone:     reqStruct.Phone,
		Email:     reqStruct.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := dao.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
	// 验证通过后删除激活码
	dao.DeleteActivationCode(reqStruct.Identifier)
}

func LoginUserService(c *gin.Context, req interface{}) {
	reqStruct, ok := req.(struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数类型错误"})
		return
	}
	if strings.TrimSpace(reqStruct.Identifier) == "" || !isValidUserInput(reqStruct.Password) {
		c.JSON(http.StatusOK, gin.H{"error": "账号和密码不能为空且不能包含空格"})
		return
	}
	var user *model.User
	var err error
	if strings.Contains(reqStruct.Identifier, "@") {
		user, err = dao.GetUserByEmail(reqStruct.Identifier)
	} else {
		user, err = dao.GetUserByPhone(reqStruct.Identifier)
	}
	if err != nil || user == nil || user.Password != reqStruct.Password {
		c.JSON(http.StatusOK, gin.H{"error": "账号或密码错误"})
		return
	}

	// 生成 AccessToken 和 RefreshToken
	accessToken, err := generateAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成access token失败"})
		return
	}

	refreshToken, err := generateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成refresh token失败"})
		return
	}

	// 将 RefreshToken 存储到 Redis 中，用于后续验证
	err = RedisClient.Set(c.Request.Context(), "refresh_token:"+string(rune(user.ID)), refreshToken, 7*24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "存储refresh token失败"})
		return
	}

	tokenResponse := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    15 * 60, // 15分钟
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"tokens":  tokenResponse,
		"user":    user,
	})
}

// 刷新 Token 服务
func RefreshTokenService(c *gin.Context, req interface{}) {
	reqStruct, ok := req.(struct {
		RefreshToken string `json:"refresh_token"`
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数类型错误"})
		return
	}

	// 验证 RefreshToken
	claims, err := validateToken(reqStruct.RefreshToken, "refresh")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token无效"})
		return
	}

	// 检查 RefreshToken 是否在 Redis 中存在
	storedToken, err := RedisClient.Get(c.Request.Context(), "refresh_token:"+string(rune(claims.UserID))).Result()
	if err != nil || storedToken != reqStruct.RefreshToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token已失效"})
		return
	}

	// 生成新的 AccessToken
	newAccessToken, err := generateAccessToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成新access token失败"})
		return
	}

	// 可选：生成新的 RefreshToken (滚动刷新)
	newRefreshToken, err := generateRefreshToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成新refresh token失败"})
		return
	}

	// 更新 Redis 中的 RefreshToken
	err = RedisClient.Set(c.Request.Context(), "refresh_token:"+string(rune(claims.UserID)), newRefreshToken, 7*24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新refresh token失败"})
		return
	}

	tokenResponse := TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    15 * 60,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "token刷新成功",
		"tokens":  tokenResponse,
	})
}

// 登出服务 - 清除 RefreshToken
func LogoutService(c *gin.Context, req interface{}) {
	reqStruct, ok := req.(struct {
		RefreshToken string `json:"refresh_token"`
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数类型错误"})
		return
	}

	// 验证 RefreshToken 并获取用户信息
	claims, err := validateToken(reqStruct.RefreshToken, "refresh")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token无效"})
		return
	}

	// 从 Redis 中删除 RefreshToken
	err = RedisClient.Del(c.Request.Context(), "refresh_token:"+string(rune(claims.UserID))).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登出失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

func UploadUserService(c *gin.Context, req interface{}) {
	reqStruct, ok := req.(struct {
		Id          int    `json:"id"`
		Username    string `json:"username"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数类型错误"})
		return
	}
	user := &model.User{}
	db := dao.GetDB()
	if err := db.First(user, reqStruct.Id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}
	if user.Password != reqStruct.Password {
		c.JSON(http.StatusOK, gin.H{"error": "密码错误"})
		return
	}
	update := map[string]any{}
	if reqStruct.Username != "" {
		update["username"] = reqStruct.Username
	}
	if reqStruct.Phone != "" {
		update["phone"] = reqStruct.Phone
	}
	if reqStruct.Email != "" {
		if !IsValidEmail(reqStruct.Email) {
			c.JSON(http.StatusOK, gin.H{"error": "邮箱格式不正确"})
			return
		}
		update["email"] = reqStruct.Email
	}
	if reqStruct.NewPassword != "" {
		update["password"] = reqStruct.NewPassword
	}
	if len(update) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": "无可更新字段"})
		return
	}
	if err := db.Model(user).Updates(update).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func GetVerificationCodeService(c *gin.Context, req interface{}) {
	reqStruct, ok := req.(struct {
		Email string `json:"email"`
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数类型错误"})
		return
	}
	if exists, _ := RedisClient.Exists(c.Request.Context(), reqStruct.Email).Result(); exists > 0 {
		c.JSON(200, "请勿频繁请求")
		return
	}
	digits := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	codeRunes := make([]rune, 6)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 6; i++ {
		idx := rand.Intn(10)
		codeRunes[i] = digits[idx]
	}
	code := string(codeRunes)
	if !IsValidEmail(reqStruct.Email) {
		c.JSON(http.StatusOK, gin.H{"error": "邮箱格式不正确"})
		return
	}
	if err := sendEmailVerificationCode(reqStruct.Email, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送邮箱验证码失败"})
		return
	}
	err := RedisClient.Set(c.Request.Context(), reqStruct.Email, code, 3*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "验证码存储失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "邮箱验证码已生成", "code": code})
}

func isValidUserInput(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}
	if strings.ContainsAny(s, " ") {
		return false
	}
	return true
}

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func checkVerificationCode(c *gin.Context, email string, code string) error {
	if correctCode, err := RedisClient.Get(c.Request.Context(), email).Result(); err != nil {
		return errors.New("验证码不存在")
	} else if correctCode == code {
		RedisClient.Del(c.Request.Context(), email)
		return nil
	}
	return errors.New("无效的验证码")
}

func sendEmailVerificationCode(email, code string) error {
	from := config.Conf.GetString("qq.from")
	password := config.Conf.GetString("qq.password")
	smtpHost := "smtp.qq.com"
	smtpPort := 587

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "您的验证码")
	m.SetBody("text/plain", "您的验证码是："+code+"，1分钟内有效。")

	d := gomail.NewDialer(smtpHost, smtpPort, from, password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
