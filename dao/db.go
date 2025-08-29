package dao

import (
	"group_ten_server/config"
	"group_ten_server/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	dsn := config.Conf.GetString("mysql.dsn")
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	log.Println("数据库连接成功")
	db = database

	// 自动迁移所有模型表结构
	err = db.AutoMigrate(
		&model.User{},
		&model.Environment{},
		&model.ActivationCode{},
		&model.Ticket{},
	)
	if err != nil {
		log.Fatalf("自动迁移表结构失败: %v", err)
	}
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}
