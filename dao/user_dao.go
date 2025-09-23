package dao

import (
	"github.com/xiuivfbc/smart_museum/model"
	"gorm.io/gorm"
)

// GetUserByPhone 根据电话查找用户
func GetUserByPhone(phone string) (*model.User, error) {
	var user model.User
	err := db.Where("phone = ?", phone).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// GetUserByEmail 根据邮箱查找用户
func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := db.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// CreateUser 创建新用户
func CreateUser(user *model.User) error {
	return db.Create(user).Error
}

// GetUserIDByEmailOrPhone 通过邮箱或电话获取用户ID
func GetUserIDByEmailOrPhone(identifier string) (int, error) {
	var user model.User
	err := db.Where("email = ? OR phone = ?", identifier, identifier).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}
