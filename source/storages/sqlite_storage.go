package storages

import (
	"database/sql"
	_ "embed"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

var _ Storage = &SQLiteStorage{}

type SQLiteStorage struct {
	dbPath        string
	musicFilePath string
	db            *sql.DB
	user          User
	filesystemMux sync.RWMutex
}

// NewSQLiteStorage initializes the storage with the provided user login info.
// Logs fatal if anything goes wrong.
func NewSQLiteStorage(dbPath string, musicFilePath string, user User) *SQLiteStorage {
	storage := &SQLiteStorage{
		dbPath:        dbPath,
		musicFilePath: musicFilePath,
		user:          user,
	}

	// Create database.
	var err error
	storage.db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		slog.Error("failed to open SQLite database", "error", err)
		return nil
	}
	if _, err := storage.db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		slog.Error("failed to enable foreign keys", "error", err)
		return nil
	}
	if _, err := storage.db.Exec(schemaSQL); err != nil {
		slog.Error("failed to create schema", "error", err)
		return nil
	}

	// Create music file directory.
	if err := os.MkdirAll(musicFilePath, 0700); err != nil {
		slog.Error("failed to create music file directory", "error", err)
		return nil
	}
	return storage
}

// ---------------- Playlist Methods ----------------

func (s *SQLiteStorage) CreatePlaylist(title string, coverBlob []byte) (Playlist, error) {
	currentTime := time.Now().UnixNano()

	p := Playlist{
		UserId:       s.user.UserId,
		PlaylistId:   currentTime,
		Title:        title,
		ModifiedDate: currentTime,
		CoverBlob:    coverBlob,
	}

	_, err := s.db.Exec(
		`INSERT OR IGNORE INTO playlists (user_id, playlist_id, title, modified_date, cover_blob)
		VALUES (?, ?, ?, ?, ?)`,
		p.UserId, p.PlaylistId, p.Title, p.ModifiedDate, p.CoverBlob,
	)
	return p, err
}

func (s *SQLiteStorage) UpdatePlaylist(p Playlist) error {
	_, err := s.db.Exec(
		`UPDATE playlists SET title = ?, modified_date = ?, cover_blob = ?
		WHERE user_id = ? AND playlist_id = ?`,
		p.Title, time.Now().UnixNano(), p.CoverBlob, s.user.UserId, p.PlaylistId,
	)
	return err
}

func (s *SQLiteStorage) DeletePlaylist(playlistID int64) error {
	_, err := s.db.Exec(
		`DELETE FROM playlists WHERE user_id = ? AND playlist_id = ?`,
		s.user.UserId, playlistID,
	)
	return err
}

func (s *SQLiteStorage) GetPlaylist(playlistID int64) (Playlist, error) {
	var p Playlist
	p.UserId = s.user.UserId

	row := s.db.QueryRow(
		`SELECT playlist_id, title, modified_date, cover_blob
		FROM playlists WHERE user_id = ? AND playlist_id = ?`,
		s.user.UserId, playlistID,
	)

	err := row.Scan(&p.PlaylistId, &p.Title, &p.ModifiedDate, &p.CoverBlob)
	return p, err
}

func (s *SQLiteStorage) GetAllSortedPlaylists() ([]Playlist, error) {
	rows, err := s.db.Query(
		`SELECT playlist_id, title, modified_date, cover_blob
		FROM playlists WHERE user_id = ?
		ORDER BY modified_date DESC`,
		s.user.UserId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []Playlist
	for rows.Next() {
		var p Playlist
		p.UserId = s.user.UserId
		if err := rows.Scan(&p.PlaylistId, &p.Title, &p.ModifiedDate, &p.CoverBlob); err != nil {
			return nil, err
		}
		playlists = append(playlists, p)
	}
	return playlists, nil
}

// ---------------- Music Methods ----------------

func (s *SQLiteStorage) CreateMusic(m Music) error {
	_, err := s.db.Exec(
		`INSERT OR IGNORE INTO music (music_id, source, title, length_seconds)
		VALUES (?, ?, ?, ?)`,
		m.MusicId, m.Source, m.Title, m.LengthSeconds,
	)
	return err
}

func (s *SQLiteStorage) UpdateMusic(m Music) error {
	_, err := s.db.Exec(
		`UPDATE music SET title = ?, length_seconds = ?
		WHERE music_id = ? AND source = ?`,
		m.Title, m.LengthSeconds, m.MusicId, m.Source,
	)
	return err
}

func (s *SQLiteStorage) DeleteMusic(musicID string, source MusicSource) error {
	_, err := s.db.Exec(`DELETE FROM music WHERE music_id = ? AND source = ?`,
		musicID, source)
	return err
}

func (s *SQLiteStorage) GetMusic(musicID string, source MusicSource) (Music, error) {
	var m Music
	row := s.db.QueryRow(
		`SELECT music_id, source, title, length_seconds
		FROM music WHERE music_id = ? AND source = ?`,
		musicID, source,
	)
	err := row.Scan(&m.MusicId, &m.Source, &m.Title, &m.LengthSeconds)
	return m, err
}

func (s *SQLiteStorage) GetAllMusic() ([]Music, error) {
	rows, err := s.db.Query(`SELECT music_id, source, title, length_seconds FROM music`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var musics []Music
	for rows.Next() {
		var m Music
		if err := rows.Scan(&m.MusicId, &m.Source, &m.Title, &m.LengthSeconds); err != nil {
			return nil, err
		}
		musics = append(musics, m)
	}
	return musics, nil
}

// ---------------- Playlist-Music Methods ----------------

func (s *SQLiteStorage) AddMusicToPlaylist(playlistID int64, musicID string, source MusicSource) error {
	currentTime := time.Now().UnixNano()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec(
		`INSERT OR IGNORE INTO playlist_music (user_id, playlist_id, music_id, source, added_at)
		VALUES (?, ?, ?, ?, ?)`,
		s.user.UserId, playlistID, musicID, source, currentTime,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Check if a row was actually added before unnecessarily updating the playlist's modified date
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		_, err = tx.Exec(
			`UPDATE playlists SET modified_date = ? 
			WHERE user_id = ? AND playlist_id = ?`,
			currentTime, s.user.UserId, playlistID,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *SQLiteStorage) RemoveMusicFromPlaylist(playlistID int64, musicID string, source MusicSource) error {
	currentTime := time.Now().UnixNano()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM playlist_music 
		WHERE user_id = ? AND playlist_id = ? AND music_id = ? AND source = ?`,
		s.user.UserId, playlistID, musicID, source,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		`UPDATE playlists SET modified_date = ? 
		WHERE user_id = ? AND playlist_id = ?`,
		currentTime, s.user.UserId, playlistID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *SQLiteStorage) GetAllSortedMusicFromPlaylist(playlistID int64) ([]Music, error) {
	rows, err := s.db.Query(
		`SELECT m.music_id, m.source, m.title, m.length_seconds
		FROM music m
		JOIN playlist_music pm ON m.music_id = pm.music_id AND m.source = pm.source
		WHERE pm.user_id = ? AND pm.playlist_id = ?
		ORDER BY pm.added_at ASC`,
		s.user.UserId, playlistID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var musics []Music
	for rows.Next() {
		var m Music
		if err := rows.Scan(&m.MusicId, &m.Source, &m.Title, &m.LengthSeconds); err != nil {
			return nil, err
		}
		musics = append(musics, m)
	}
	return musics, nil
}

// ---------------- Filesystem Methods ----------------

func (s *SQLiteStorage) getMusicFilePath(music Music) string {
	return filepath.Join(s.musicFilePath, fmt.Sprintf("%v_%v.mp3", music.Source, music.MusicId))
}

func (s *SQLiteStorage) CreateOrUpdateMusicFile(music Music, content io.Reader) error {
	s.filesystemMux.Lock()
	defer s.filesystemMux.Unlock()

	file, err := os.OpenFile(s.getMusicFilePath(music), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, content)
	return err
}

func (s *SQLiteStorage) RemoveMusicFile(music Music) error {
	s.filesystemMux.Lock()
	defer s.filesystemMux.Unlock()
	return os.Remove(s.getMusicFilePath(music))
}

func (s *SQLiteStorage) GetMusicFile(music Music) (io.ReadCloser, error) {
	s.filesystemMux.RLock()
	defer s.filesystemMux.RUnlock()
	return os.Open(s.getMusicFilePath(music))
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
