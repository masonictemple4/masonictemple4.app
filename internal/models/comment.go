package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	PostID uint `gorm:"column:postid;index;" json:"postid"`
	UserID uint `gorm:"column:userid;index;" json:"userid"`
	User   User `json:"user"`
	Text   uint `gorm:"column:text" json:"text"`
	// TODO: add likes.
}
