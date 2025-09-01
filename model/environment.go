package model

import "group_ten_server/config"

// Environment 表示环境数据模型
// 包含编号、温度、湿度、状态、信息、灯光、火焰等字段
// 可用于数据库映射和前后端交互

type Environment struct {
	ID          int     `json:"id" gorm:"primaryKey"` // 编号
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Gas         int     `json:"gas"`
	AutoFan     int     `json:"auto_fan"`
	FanLevel    int     `json:"fan_level"`
	AutoLight   int     `json:"auto_light"`
	Light       int     `json:"light"`
	Human       int     `json:"human"` // 通常表示是否检测到人 (0: 否, 1: 是)
	Dark        int     `json:"dark"`  // 通常表示环境是否黑暗 (0: 否, 1: 是)
}

// TableName 动态表名方法
func (Environment) TableName(roomName string) string {
	return config.AppConfigInstance.RoomMapping[roomName]
}
