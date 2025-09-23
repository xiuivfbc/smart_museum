package dao

import (
	"time"

	"github.com/xiuivfbc/smart_museum/model"
)

func InsertNewDailyEntryCount(v int) error {
	now := time.Now()
	yesterday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
	temp := model.Daily{
		Time: yesterday,
		Num:  v,
	}
	// 插入新行
	if err := db.Create(&temp).Error; err != nil {
		return err
	}
	return nil
}

// GetAllEntryCounts 查询所有每日人数统计
func GetAllEntryCounts() ([]model.Daily, error) {
	var counts []model.Daily
	err := db.Order("time desc").Find(&counts).Error
	return counts, err
}
