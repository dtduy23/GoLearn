package song

type SongUploadRequest struct {
	Title     string   `form:"title" json:"title" binding:"required,max=255"`
	AlbumID   string   `form:"album_id" json:"album_id"`     // optional
	ArtistIDs []string `form:"artist_ids" json:"artist_ids"` // optional
	GenreIDs  []string `form:"genre_ids" json:"genre_ids"`   // optional
}

type SongUploadResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	Message  string `json:"message"`
}
