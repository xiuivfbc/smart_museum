package dao

import (
	"group_ten_server/model"
)

func CreateEnvironmentToTable(table string, env *model.Environment) error {
	return db.Table(table).Create(env).Error
}

func GetEnvironmentByNameFromTable(table string) (*model.Environment, error) {
	var env model.Environment
	err := db.Table(table).
		Order("id DESC"). // 按ID倒序
		First(&env).Error // 取第一条（即ID最大的）
	return &env, err
}

func GetAllEnvironmentsFromTable(table string) ([]model.Environment, error) {
	var envs []model.Environment
	err := db.Table(table).Find(&envs).Error
	return envs, err
}

func DeleteEnvironmentByNameFromTable(table string) error {
	return db.Table(table).Delete(&model.Environment{}).Error
}
