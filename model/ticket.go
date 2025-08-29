package model

// User 表示用户信息的数据模型
// 包含ID、用户名、密码、邮箱、角色、创建时间等字段，适合JWT身份验证扩展
import "time"

type Ticket struct {
	ID   int       `json:"id" gorm:"primaryKey"` // 用户唯一ID
	Path string    `json:"path"`                 // 路径
	Time time.Time `json:"time"`                 // 更新时间
}
