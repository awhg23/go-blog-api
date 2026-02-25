package model

import "time"

// User 对应数据库里的 users 表
type User struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"`
	Username       string    `gorm:"type:varchar(50);not null;unique"`
	PasswordDigest string    `gorm:"type:varchar(255);not null"`
	CreatedAt      time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}
