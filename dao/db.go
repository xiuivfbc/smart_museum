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

	// 自动迁移模型表结构
	if err := db.AutoMigrate(
		&model.User{},
		&model.Environment{},
		&model.ActivationCode{},
		&model.Ticket{},
	); err != nil {
		log.Fatalf("自动迁移表失败: %v", err)
	}

	// 动态创建表
	if err := AutoCreateTables(); err != nil {
		log.Fatalf("动态创建表失败: %v", err)
	}
}

func AutoCreateTables() error {
	for roomName, tableName := range config.AppConfigInstance.RoomMapping {
		// 创建动态表
		err := db.Table(tableName).AutoMigrate(&model.Environment{})
		if err != nil {
			log.Printf("Failed to create table %s for room %s: %v", tableName, roomName, err)
			return err
		}
	}
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}
