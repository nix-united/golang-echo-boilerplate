package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title   string `json:"title" gorm:"type:text"`
	Content string `json:"content" gorm:"type:text"`
	UserID  uint
	User    User `gorm:"foreignkey:UserID"`
}
