package dtos

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type GoogleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type AuthResponse struct {
	JWT         string    `json:"jwt_token"`
	AccessToken string    `json:"access_token,omitempty"`
	ExpiresAt   time.Time `json:"expires_at"`
	User        UserInfo  `json:"user"`
}

type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Claims struct {
	UserID   string `json:"user_id"`
	GoogleID string `json:"google_id"`
	FolderID string `json:"folder_id"`
	jwt.RegisteredClaims
}
