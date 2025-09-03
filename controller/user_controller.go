package controller

import (
	"errors"
	"group_ten_server/config"
	"group_ten_server/dao"
	"group_ten_server/model"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
)

var jwtKey []byte

var RedisClient *redis.Client

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
func RegisterUser(c *gin.Context) {
	var req struct {
		Username         string `json:"username"`
		Password         string `json:"password"`
		Role             string `json:"role"`
		Phone            string `json:"phone"`
		Email            string `json:"email"`
		Identifier       string `json:"identifier"`        // admin验证码
		VerificationCode string `json:"verification_code"` // 验证码
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if !isValidUserInput(req.Password) {
		c.JSON(http.StatusOK, gin.H{"error": "密码不能为空且不能包含空格"})
		return
	}
	if strings.TrimSpace(req.Phone) == "" && strings.TrimSpace(req.Email) == "" {
		c.JSON(http.StatusOK, gin.H{"error": "电话和邮箱至少填写一个"})
		return
	}
	if req.Email != "" && !IsValidEmail(req.Email) {
		c.JSON(http.StatusOK, gin.H{"error": "邮箱格式不正确"})
		return
	}
	// 检查电话唯一
	if req.Phone != "" {
		var user model.User
		db := dao.GetDB()
		if err := db.Where("phone = ?", req.Phone).First(&user).Error; err == nil {
			c.JSON(http.StatusOK, gin.H{"error": "电话已被注册"})
			return
		}
	}
	// 检查邮箱唯一
	if req.Email != "" {
		var user model.User
		db := dao.GetDB()
		if err := db.Where("email = ?", req.Email).First(&user).Error; err == nil {
			c.JSON(http.StatusOK, gin.H{"error": "邮箱已被注册"})
			return
		}
	}
	if ok, _ := dao.GetActivationCode(req.Identifier); req.Role == "admin" && !ok {
		c.JSON(http.StatusOK, gin.H{"error": "无效的管理员激活码"})
		return
	} else {
		dao.DeleteActivationCode(req.Identifier)
	}

	if err := checkVerificationCode(c, req.Email, req.VerificationCode); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
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
func LoginUser(c *gin.Context) {
	var req struct {
		Identifier string `json:"identifier"` // 电话或邮箱
		Password   string `json:"password"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if strings.TrimSpace(req.Identifier) == "" || !isValidUserInput(req.Password) {
		c.JSON(http.StatusOK, gin.H{"error": "账号和密码不能为空且不能包含空格"})
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
		c.JSON(http.StatusOK, gin.H{"error": "账号或密码错误"})
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
	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "token": tokenString, "id": user.ID, "userform": user})
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
	// 查询用户qfda
	user := &model.User{}
	db := dao.GetDB()
	if err := db.First(user, req.Id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在"})
		return
	}
	// 校验密码
	if user.Password != req.Password {
		c.JSON(http.StatusOK, gin.H{"error": "密码错误"})
		return
	}
	// 更新信息
	update := map[string]any{}
	if req.Username != "" {
		update["username"] = req.Username
	}
	if req.Phone != "" {
		update["phone"] = req.Phone
	}
	if req.Email != "" {
		if !IsValidEmail(req.Email) {
			c.JSON(http.StatusOK, gin.H{"error": "邮箱格式不正确"})
			return
		}
		update["email"] = req.Email
	}
	if req.NewPassword != "" {
		update["password"] = req.NewPassword
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

func GetVerificationCode(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	// 检查是否发送太频繁 (Redis里查这个号有没有记录)
	if exists, _ := RedisClient.Exists(c.Request.Context(), req.Email).Result(); exists > 0 {
		c.JSON(200, "请勿频繁请求")
		return
	}
	// 生成6位随机数字验证码
	digits := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	codeRunes := make([]rune, 6)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 6; i++ {
		idx := rand.Intn(10)
		codeRunes[i] = digits[idx]
	}
	code := string(codeRunes)

	// 邮箱
	if !IsValidEmail(req.Email) {
		c.JSON(http.StatusOK, gin.H{"error": "邮箱格式不正确"})
		return
	}
	if err := sendEmailVerificationCode(req.Email, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送邮箱验证码失败"})
		return
	}
	// 将验证码存入Redis
	err := RedisClient.Set(c.Request.Context(), req.Email, code, 3*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "验证码存储失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "邮箱验证码已生成", "code": code})
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
