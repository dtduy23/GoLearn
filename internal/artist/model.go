package artist

import (
	"time"

	"github.com/google/uuid"
)

// Core artist info
type Artist struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Profile info for artist
type ArtistProfile struct {
	ArtistID         uuid.UUID `json:"artist_id"`
	Bio              string    `json:"bio"`
	ImageURL         string    `json:"image_url"`
	MonthlyListeners int       `json:"monthly_listeners"`
	Country          string    `json:"country"`
	Verified         bool      `json:"verified"`
}
