package models

import "gorm.io/gorm"

type Providers string

const GOOGLE Providers = "google"

type OAuthProviders struct {
	gorm.Model
	UserID   uint      `json:"user_id"`
	Token    string    `json:"token"`
	Provider Providers `json:"provider"`
}
