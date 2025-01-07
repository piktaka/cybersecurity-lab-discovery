package models

import (
	"lablabee.com/cybersecurity-discovery1/hping-plateform/database"
)

type Post struct {
	ID       uint      `gorm:"primaryKey;autoIncrement"`
	Content  string    `gorm:"type:text;not null"`
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

func CreatePost(content string) (*Post, error) {
	post := &Post{Content: content}
	result := database.DB.Create(post)
	return post, result.Error
}

func AddComment(comment Comment) (*Post, error) {
	post := &Post{}

	result := database.DB.Where(comment.PostID).First(post)
	if err := result.Error; err != nil {
		return nil, err
	}
	post.Comments = append(post.Comments, comment)
	result = database.DB.Save(post)
	if err := result.Error; err != nil {
		return nil, err
	}
	return post, nil
}

func GetAllPosts() ([]Post, error) {
	posts := make([]Post, 0, 10)
	err := database.DB.Preload("Comments").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}
