package models

import "lablabee.com/cybersecurity-discovery1/hping-plateform/database"

type Comment struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	PostID  uint   `gorm:"not null"`
	Content string `gorm:"type:text;not null"`
}

func CreateComment(postID uint, content string) (*Comment, error) {
	comment := Comment{PostID: postID, Content: content}
	result := database.DB.Create(&comment)
	return &comment, result.Error
}
