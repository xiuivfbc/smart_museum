package dao

import (
	"github.com/xiuivfbc/smart_museum/model"
	"gorm.io/gorm"
)

// GetActivationCode 检查激活码是否存在
func GetActivationCode(code string) (bool, error) {
	var ac model.ActivationCode
	err := db.First(&ac, "code = ?", code).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteActivationCode 删除已使用的激活码
func DeleteActivationCode(code string) error {
	return db.Delete(&model.ActivationCode{}, "code = ?", code).Error
}

// EnsureActivationCodes 保证激活码数量为5，不足则补充
func EnsureActivationCodes() error {
	count, err := CountActivationCodes()
	if err != nil {
		return err
	}
	if count >= 5 {
		return nil
	}
	for i := count; i < 5; i++ {
		newCode := model.GenerateActivationCode()
		// 防止重复
		exist, err := GetActivationCode(newCode)
		if err != nil {
			return err
		}
		if exist {
			i--
			continue
		}
		err = AddActivationCode()
		if err != nil {
			return err
		}
	}
	return nil
}

// AddActivationCode 添加新激活码
func AddActivationCode() error {
	code := model.GenerateActivationCode()
	return db.Create(&model.ActivationCode{Code: code}).Error
}

// CountActivationCodes 获取当前激活码数量
func CountActivationCodes() (int64, error) {
	var count int64
	err := db.Model(&model.ActivationCode{}).Count(&count).Error
	return count, err
}
