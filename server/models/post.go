package models

import "github.com/jinzhu/gorm"

type Post struct {
	gorm.Model
	Title   string `json:"title" gorm:"type:text"`
	Content string `json:"content" gorm:"type:text"`
	UserId  uint
	User    User   `gorm:"foreignkey:UserId"`
}
