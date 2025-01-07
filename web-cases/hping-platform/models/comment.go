package models

import (
	"time"

	"lablabee.com/cybersecurity-discovery1/hping-plateform/database"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	PostID    uint      `gorm:"not null"`
	Content   string    `gorm:"type:text;not null"`
	Timestamp time.Time `gorm:"not null"`
}

func CreateComment(postID uint, content string, timestamp time.Time) (*Comment, error) {
	comment := Comment{PostID: postID, Content: content, Timestamp: timestamp}
	result := database.DB.Create(&comment)
	return &comment, result.Error
}
