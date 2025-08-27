package dao

import (
	"group_ten_server/model"

	"gorm.io/gorm"
)

// GetUserByUsername 根据用户名查找用户
func GetUserByUsername(username string) (*model.User, error) {
	db := GetDB()
	var user model.User
	err := db.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// CreateUser 创建新用户
func CreateUser(user *model.User) error {
	db := GetDB()
	return db.Create(user).Error
}
