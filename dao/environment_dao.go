package dao

import (
	"group_ten_server/model"
	"log"
)

// InitEnvironmentTable 自动迁移 Environment 表结构
func InitEnvironmentTable() {
	if db == nil {
		log.Fatal("数据库未初始化")
	}
	if err := db.AutoMigrate(&model.Environment{}); err != nil {
		log.Fatalf("自动迁移 Environment 表失败: %v", err)
	}
}

// CreateEnvironment 新增环境数据
func CreateEnvironment(env *model.Environment) error {
	return db.Create(env).Error
}

// GetEnvironmentByName 根据名称查询环境数据
func GetEnvironmentByName(name string) (*model.Environment, error) {
	var env model.Environment
	err := db.Where("name = ?", name).First(&env).Error
	return &env, err
}

// GetAllEnvironments 查询所有环境数据
func GetAllEnvironments() ([]model.Environment, error) {
	var envs []model.Environment
	err := db.Find(&envs).Error
	return envs, err
}

// UpdateEnvironmentByName 根据name更新环境数据
func UpdateEnvironmentByName(name string, update *model.Environment) error {
	return db.Model(&model.Environment{}).Where("name = ?", name).Updates(update).Error
}

// DeleteEnvironmentByName 根据name删除环境数据
func DeleteEnvironmentByName(name string) error {
	return db.Where("name = ?", name).Delete(&model.Environment{}).Error
}
