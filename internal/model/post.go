package model

import "time"

// Post 对应数据库里的 posts 表
type Post struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID    int64     `gorm:"not null;index"` //加上普通索引加速查询
	User      User      `gorm:"foreignKey:UserID"`
	Title     string    `gorm:"type:varchar(100);not null"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}
