package model

import (
	"time"

	"github.com/google/uuid"
)

// Core authentication
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Profile info
type UserProfile struct {
	ShowProfile bool      `json:"show_profile"`
	UserID      uuid.UUID `json:"user_id"`
	FullName    string    `json:"full_name"`
	AvatarURL   string    `json:"avatar_url"`
	Sex         string    `json:"sex"`
	Birthday    time.Time `json:"birthday"`
	Country     string    `json:"country"`
}
