package song

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// GetByID fetches a song by ID with album, artists, and genres
func (r *Repository) GetByID(ctx context.Context, id string) (*Song, error) {
	// 1. Get song with album info
	song, err := r.getSongWithAlbum(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. Get artists for this song
	artists, err := r.getArtistsBySongID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching artists: %w", err)
	}
	song.Artists = artists

	// 3. Get genres for this song
	genres, err := r.getGenresBySongID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching genres: %w", err)
	}
	song.Genres = genres

	return song, nil
}

// getSongWithAlbum fetches song with album data using LEFT JOIN
func (r *Repository) getSongWithAlbum(ctx context.Context, id string) (*Song, error) {
	query := `
		SELECT 
			s.id, s.title, s.duration, s.file_url, s.play_count, 
			COALESCE(s.track_number, 0), s.created_at,
			a.id, a.title, a.cover_url
		FROM songs s
		LEFT JOIN albums a ON s.album_id = a.id
		WHERE s.id = $1
	`

	var song Song
	var albumID, albumTitle, albumCoverURL *string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&song.ID,
		&song.Title,
		&song.Duration,
		&song.FileURL,
		&song.PlayCount,
		&song.TrackNumber,
		&song.CreatedAt,
		&albumID,
		&albumTitle,
		&albumCoverURL,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("song not found")
		}
		return nil, fmt.Errorf("error querying song: %w", err)
	}

	// Set album if exists
	if albumID != nil {
		song.Album = &Album{
			ID:       *albumID,
			Title:    stringOrEmpty(albumTitle),
			CoverURL: stringOrEmpty(albumCoverURL),
		}
	}

	return &song, nil
}

// getArtistsBySongID fetches all artists for a song
func (r *Repository) getArtistsBySongID(ctx context.Context, songID string) ([]SongArtist, error) {
	query := `
		SELECT a.id, a.name, COALESCE(a.image_url, ''), sa.is_primary
		FROM artists a
		INNER JOIN song_artists sa ON a.id = sa.artist_id
		WHERE sa.song_id = $1
		ORDER BY sa.is_primary DESC, a.name ASC
	`

	rows, err := r.db.Query(ctx, query, songID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artists []SongArtist
	for rows.Next() {
		var artist SongArtist
		if err := rows.Scan(&artist.ID, &artist.Name, &artist.ImageURL, &artist.IsPrimary); err != nil {
			return nil, err
		}
		artists = append(artists, artist)
	}

	return artists, rows.Err()
}

// getGenresBySongID fetches all genres for a song
func (r *Repository) getGenresBySongID(ctx context.Context, songID string) ([]Genre, error) {
	query := `
		SELECT g.id, g.name
		FROM genres g
		INNER JOIN song_genres sg ON g.id = sg.genre_id
		WHERE sg.song_id = $1
		ORDER BY g.name ASC
	`

	rows, err := r.db.Query(ctx, query, songID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []Genre
	for rows.Next() {
		var genre Genre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}

	return genres, rows.Err()
}

// CreateSong tạo song với related data (album, artists, genres)
// Sử dụng transaction để đảm bảo data consistency
func (r *Repository) CreateSong(ctx context.Context, input CreateSongInput) error {
	// Bắt đầu transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Rollback nếu có lỗi

	// 1. Insert song (với album_id nếu có)
	songQuery := `
		INSERT INTO songs (id, title, duration, file_url, play_count, track_number, album_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.Exec(ctx, songQuery,
		input.Song.ID,
		input.Song.Title,
		input.Song.Duration,
		input.Song.FileURL,
		input.Song.PlayCount,
		input.Song.TrackNumber,
		input.AlbumID, // có thể nil
		input.Song.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("error inserting song: %w", err)
	}

	// 2. Insert song_artists (nếu có)
	if len(input.ArtistIDs) > 0 {
		artistQuery := `
			INSERT INTO song_artists (song_id, artist_id, is_primary)
			VALUES ($1, $2, $3)
		`
		for i, artistID := range input.ArtistIDs {
			isPrimary := (i == 0) // Artist đầu tiên là primary
			_, err = tx.Exec(ctx, artistQuery, input.Song.ID, artistID, isPrimary)
			if err != nil {
				return fmt.Errorf("error inserting song_artist: %w", err)
			}
		}
	}

	// 3. Insert song_genres (nếu có)
	if len(input.GenreIDs) > 0 {
		genreQuery := `
			INSERT INTO song_genres (song_id, genre_id)
			VALUES ($1, $2)
		`
		for _, genreID := range input.GenreIDs {
			_, err = tx.Exec(ctx, genreQuery, input.Song.ID, genreID)
			if err != nil {
				return fmt.Errorf("error inserting song_genre: %w", err)
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// Helper function
func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
