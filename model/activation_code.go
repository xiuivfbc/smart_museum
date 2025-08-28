package model

import (
	"math/rand"
	"time"
)

// ActivationCode 表示激活码
// Code: 激活码内容，定长10

type ActivationCode struct {
	Code string `gorm:"type:char(10);primaryKey" json:"code"`
}

// GenerateActivationCode 生成一个10位随机激活码
func GenerateActivationCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 10)

	// 创建新的随机源
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}

	return string(result)
}
