package model

import "time"

type Comment struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    int64     `gorm:"not null;index" json:"post_id"` // 关联的文章ID(索引加速查询)
	UserID    int64     `gorm:"not null;index" json:"user_id"` // 评论的作者ID(索引加速查询)
	User      User      `gorm:"foreignKey:UserID" json:"user"` //预加载关联：评论者信息
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
}
