package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	GoogleID     string         `json:"google_id" gorm:"uniqueIndex;not null"`
	FolderID     string         `json:"folder_id"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null"`
	RefreshToken string         `json:"refresh_token" gorm:"uniqueIndex;not null"`
	AccessToken  string         `json:"access_token" gorm:"uniqueIndex;not null"`
	Expiry       time.Time      `json:"expiry" gorm:"not null"`
	Name         string         `json:"name"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
