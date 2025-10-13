package storage

import (
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

var _ StorageManager = &SQLiteStorageManager{}

type SQLiteStorageManager struct {
	db   *sql.DB
	user User
}

// NewSQLiteStorageManager initializes the manager for a specific user.
// Logs fatal if anything goes wrong.
func NewSQLiteStorageManager(user User) *SQLiteStorageManager {
	db, err := sql.Open("sqlite", "local.db")
	if err != nil {
		log.Fatalf("failed to open SQLite database: %v", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("failed to enable foreign keys: %v", err)
	}

	return &SQLiteStorageManager{
		db:   db,
		user: user,
	}
}

// ---------------- Playlist Methods ----------------

func (s *SQLiteStorageManager) CreatePlaylist(p Playlist) (Playlist, error) {
	currentTime := time.Now().UnixNano()
	p.PlaylistID = currentTime

	_, err := s.db.Exec(
		`INSERT INTO playlists (user_id, playlist_id, title, modified_date, cover_blob)
		 VALUES (?, ?, ?, ?, ?)`,
		s.user.UserID, p.PlaylistID, p.Title, currentTime, p.CoverBlob,
	)
	return p, err
}

func (s *SQLiteStorageManager) UpdatePlaylist(p Playlist) error {
	_, err := s.db.Exec(
		`UPDATE playlists SET title = ?, modified_date = ?, cover_blob = ?
		 WHERE user_id = ? AND playlist_id = ?`,
		p.Title, time.Now().UnixNano(), p.CoverBlob, s.user.UserID, p.PlaylistID,
	)
	return err
}

func (s *SQLiteStorageManager) DeletePlaylist(playlistID int64) error {
	_, err := s.db.Exec(
		`DELETE FROM playlists WHERE user_id = ? AND playlist_id = ?`,
		s.user.UserID, playlistID,
	)
	return err
}

func (s *SQLiteStorageManager) GetPlaylist(playlistID int64) (Playlist, error) {
	var p Playlist
	row := s.db.QueryRow(
		`SELECT playlist_id, title, cover_blob
		 FROM playlists WHERE user_id = ? AND playlist_id = ?`,
		s.user.UserID, playlistID,
	)
	err := row.Scan(&p.PlaylistID, &p.Title, &p.CoverBlob)
	return p, err
}

func (s *SQLiteStorageManager) ListPlaylists() ([]Playlist, error) {
	rows, err := s.db.Query(
		`SELECT playlist_id, title, cover_blob
		 FROM playlists WHERE user_id = ?`,
		s.user.UserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []Playlist
	for rows.Next() {
		var p Playlist
		if err := rows.Scan(&p.PlaylistID, &p.Title, &p.CoverBlob); err != nil {
			return nil, err
		}
		playlists = append(playlists, p)
	}
	return playlists, nil
}

// ---------------- Music Methods ----------------

func (s *SQLiteStorageManager) CreateMusic(m Music) error {
	_, err := s.db.Exec(
		`INSERT INTO music (music_id, source, title, length_seconds)
		 VALUES (?, ?, ?, ?)`,
		m.MusicID, m.Source, m.Title, m.LengthSeconds,
	)
	return err
}

func (s *SQLiteStorageManager) UpdateMusic(m Music) error {
	_, err := s.db.Exec(
		`UPDATE music SET title = ?, length_seconds = ?
		 WHERE music_id = ? AND source = ?`,
		m.Title, m.LengthSeconds, m.MusicID, m.Source,
	)
	return err
}

func (s *SQLiteStorageManager) DeleteMusic(musicID string, source MusicSource) error {
	_, err := s.db.Exec(`DELETE FROM music WHERE music_id = ? AND source = ?`,
		musicID, source)
	return err
}

func (s *SQLiteStorageManager) GetMusic(musicID string, source MusicSource) (Music, error) {
	var m Music
	row := s.db.QueryRow(
		`SELECT music_id, source, title, length_seconds
		 FROM music WHERE music_id = ? AND source = ?`,
		musicID, source,
	)
	err := row.Scan(&m.MusicID, &m.Source, &m.Title, &m.LengthSeconds)
	return m, err
}

func (s *SQLiteStorageManager) ListMusic() ([]Music, error) {
	rows, err := s.db.Query(`SELECT music_id, source, title, length_seconds FROM music`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var musics []Music
	for rows.Next() {
		var m Music
		if err := rows.Scan(&m.MusicID, &m.Source, &m.Title, &m.LengthSeconds); err != nil {
			return nil, err
		}
		musics = append(musics, m)
	}
	return musics, nil
}

// ---------------- Playlist-Music Methods ----------------

func (s *SQLiteStorageManager) AddMusicToPlaylist(playlistID int64, musicID string, source MusicSource) error {
	currentTime := time.Now().UnixNano()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Insert the new music association
	_, err = tx.Exec(
		`INSERT INTO playlist_music (user_id, playlist_id, music_id, source, added_at)
		 VALUES (?, ?, ?, ?, ?)`,
		s.user.UserID, playlistID, musicID, source, currentTime,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update the playlist's modified_date
	_, err = tx.Exec(
		`UPDATE playlists SET modified_date = ? 
		 WHERE user_id = ? AND playlist_id = ?`,
		currentTime, s.user.UserID, playlistID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *SQLiteStorageManager) RemoveMusicFromPlaylist(playlistID int64, musicID string, source MusicSource) error {
	currentTime := time.Now().UnixNano()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Delete the association
	_, err = tx.Exec(
		`DELETE FROM playlist_music 
		 WHERE user_id = ? AND playlist_id = ? AND music_id = ? AND source = ?`,
		s.user.UserID, playlistID, musicID, source,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update modified_date only if playlist exists (prevent foreign key issue)
	_, err = tx.Exec(
		`UPDATE playlists SET modified_date = ? 
		 WHERE user_id = ? AND playlist_id = ?`,
		currentTime, s.user.UserID, playlistID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *SQLiteStorageManager) ListMusicInPlaylist(playlistID int64) ([]Music, error) {
	rows, err := s.db.Query(
		`SELECT m.music_id, m.source, m.title, m.length_seconds
		 FROM music m
		 JOIN playlist_music pm ON m.music_id = pm.music_id AND m.source = pm.source
		 WHERE pm.user_id = ? AND pm.playlist_id = ?`,
		s.user.UserID, playlistID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var musics []Music
	for rows.Next() {
		var m Music
		if err := rows.Scan(&m.MusicID, &m.Source, &m.Title, &m.LengthSeconds); err != nil {
			return nil, err
		}
		musics = append(musics, m)
	}
	return musics, nil
}

func (s *SQLiteStorageManager) Close() error {
	return s.db.Close()
}
