package song

import "time"

// Artist represents basic artist info for a song
type SongArtist struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ImageURL  string `json:"image_url,omitempty"`
	IsPrimary bool   `json:"is_primary"`
}

// Album represents basic album info for a song
type Album struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	CoverURL string `json:"cover_url,omitempty"`
}

// Genre represents a music genre
type Genre struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Song represents a song with all related data
type Song struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Duration    int       `json:"duration"`
	FileURL     string    `json:"file_url"`
	PlayCount   int       `json:"play_count"`
	TrackNumber int       `json:"track_number,omitempty"`
	CreatedAt   time.Time `json:"created_at"`

	// Related data (populated via JOINs)
	Album   *Album       `json:"album,omitempty"`
	Artists []SongArtist `json:"artists,omitempty"`
	Genres  []Genre      `json:"genres,omitempty"`
}

type CreateSongInput struct {
	Song      Song
	AlbumID   *string  // optional album ID
	ArtistIDs []string // list of artist IDs (first one is primary)
	GenreIDs  []string // list of genre IDs
}
