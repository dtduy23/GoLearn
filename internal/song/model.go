package song

type Song struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Album     string `json:"album"`
	Genre     string `json:"genre"`
	Duration  int    `json:"duration"`
	PlayCount int    `json:"play_count"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
