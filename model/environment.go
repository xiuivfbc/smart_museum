package model

// Environment 表示环境数据模型
// 包含编号、温度、湿度、状态、信息、灯光、火焰等字段
// 可用于数据库映射和前后端交互

type Environment struct {
	ID       int     `json:"id" gorm:"primaryKey"` // 编号
	Name     string  `json:"name" gorm:"unique"`   // 环境名称，唯一
	Temp     float64 `json:"temp"`                 // 温度
	Humidity float64 `json:"humidity"`             // 湿度
	Status   string  `json:"status"`               // 状态
	Info     string  `json:"info"`                 // 信息
	Light    float64 `json:"light"`                // 灯光
	Flame    bool    `json:"flame"`                // 火焰（true=有火焰，false=无火焰）
}
