package dao

import (
	"group_ten_server/config"
	"group_ten_server/model"
	"time"
)

func InsertNewDailyEntryCount() error {
	now := time.Now()
	yesterday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
	temp := model.Daily{
		Time: yesterday,
		Num:  config.EntryNum,
	}
	// 插入新行
	if err := db.Create(&temp).Error; err != nil {
		return err
	}
	// 重置全局计数器
	config.EntryNum = 0
	return nil
}

// GetAllEntryCounts 查询所有每日人数统计
func GetAllEntryCounts() ([]model.Daily, error) {
	var counts []model.Daily
	err := db.Order("time desc").Find(&counts).Error
	return counts, err
}
