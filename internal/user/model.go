package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Core authentication
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
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

func (p *User) Validate() error {
	if p.Username == "" {
		return errors.New("username required")
	}
	if p.Password == "" {
		return errors.New("password required")
	}
	if p.Email == "" {
		return errors.New("email required")
	}
	return nil
}
