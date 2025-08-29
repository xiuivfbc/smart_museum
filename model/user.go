package model

// User 表示用户信息的数据模型
// 包含ID、用户名、密码、邮箱、角色、创建时间等字段，适合JWT身份验证扩展
import "time"

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`     // 用户唯一ID
	Username  string    `json:"username" `                // 用户名
	Password  string    `json:"password"`                 // 密码（加密存储）
	Role      string    `json:"role"`                     // 角色（如admin/user等）
	Phone     string    `json:"phone" gorm:"uniqueIndex"` // 电话
	Email     string    `json:"email" gorm:"uniqueIndex"` // 邮箱
	CreatedAt time.Time `json:"created_at"`               // 创建时间
	UpdatedAt time.Time `json:"updated_at"`               // 更新时间
}
