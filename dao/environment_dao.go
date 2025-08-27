package dao

import (
	"group_ten_server/model"
	"log"
)

// InitEnvironmentTable 自动迁移 Environment 表结构
func InitEnvironmentTable() {
	db := GetDB()
	if db == nil {
		log.Fatal("数据库未初始化")
	}
	if err := db.AutoMigrate(&model.Environment{}); err != nil {
		log.Fatalf("自动迁移 Environment 表失败: %v", err)
	}
}

// CreateEnvironment 新增环境数据
func CreateEnvironment(env *model.Environment) error {
	db := GetDB()
	return db.Create(env).Error
}

// GetEnvironmentByID 根据编号查询环境数据
func GetEnvironmentByID(id int) (*model.Environment, error) {
	db := GetDB()
	var env model.Environment
	err := db.First(&env, id).Error
	return &env, err
}

// GetAllEnvironments 查询所有环境数据
func GetAllEnvironments() ([]model.Environment, error) {
	db := GetDB()
	var envs []model.Environment
	err := db.Find(&envs).Error
	return envs, err
}

// UpdateEnvironment 更新环境数据
func UpdateEnvironment(env *model.Environment) error {
	db := GetDB()
	return db.Save(env).Error
}

// DeleteEnvironment 删除环境数据
func DeleteEnvironment(id int) error {
	db := GetDB()
	return db.Delete(&model.Environment{}, id).Error
}
