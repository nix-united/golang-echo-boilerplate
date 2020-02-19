package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name string  `json:"name" gorm:"type:varchar(200);"`
	Password string `json:"password" gorm:"type:varchar(200);"`
	Post []Post
}
