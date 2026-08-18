package storages

import (
	"database/sql"
	_ "embed"
	"errors"
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

var _ Storage = (*SQLiteStorage)(nil)

type SQLiteStorage struct {
	musicFilePath string
	db            *sql.DB

	cachedProfile *UserProfile
	userMutex     sync.Mutex
}

func NewSQLiteStorage(dbPath string, musicFilePath string) *SQLiteStorage {
	storage := &SQLiteStorage{
		musicFilePath: musicFilePath,
	}

	var err error
	storage.db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		slog.Error("failed to open SQLite database", "error", err)
		return nil
	}
	if _, err := storage.db.Exec(schemaSQL); err != nil {
		slog.Error("failed to create schema", "error", err)
		return nil
	}

	if err := os.MkdirAll(musicFilePath, 0700); err != nil {
		slog.Error("failed to create music file directory", "error", err)
		return nil
	}
	return storage
}

// ---------------- User Storer Methods ----------------

func (s *SQLiteStorage) PutUser(profile UserProfile) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Enforce the "at most 1 row" rule by wiping the table before insert
	if _, err := tx.Exec(`DELETE FROM user_profile`); err != nil {
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO user_profile (
			user_id, username, language, registration_date, token
		) VALUES (?, ?, ?, ?, ?)`,
		profile.UserId, profile.Username, profile.Language, profile.RegistrationDate, profile.Token,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Invalidate the cache. GetUser() will lazily fetch the single source of truth.
	s.userMutex.Lock()
	s.cachedProfile = nil
	s.userMutex.Unlock()

	return nil
}

func (s *SQLiteStorage) GetUser() (UserProfile, error) {
	s.userMutex.Lock()
	defer s.userMutex.Unlock()

	// 1. Return from cache if it exists
	if s.cachedProfile != nil {
		return *s.cachedProfile, nil
	}

	// 2. Otherwise, lazily fetch from database
	var p UserProfile
	query := `
		SELECT 
			user_id, username, language, registration_date, token
		FROM user_profile 
		LIMIT 1
	`

	err := s.db.QueryRow(query).Scan(&p.UserId, &p.Username, &p.Language, &p.RegistrationDate, &p.Token)
	if err != nil {
		return UserProfile{}, err
	}

	// 3. Populate cache
	s.cachedProfile = &p

	return p, nil
}

func (s *SQLiteStorage) DeleteUser() error {
	_, err := s.db.Exec(`DELETE FROM user_profile`)
	if err != nil {
		return err
	}

	// Invalidate cache
	s.userMutex.Lock()
	s.cachedProfile = nil
	s.userMutex.Unlock()

	return nil
}

// ---------------- Playlist Methods ----------------

func (s *SQLiteStorage) PutPlaylist(p Playlist) (Playlist, error) {
	user, err := s.GetUser()
	if err != nil {
		return Playlist{}, fmt.Errorf("failed to get user: %v", err)
	}

	currentTime := time.Now().UnixNano()
	p.ModifiedDate = currentTime
	p.UserId = user.UserId // Enforce the user ID from context

	if p.PlaylistId == 0 {
		p.PlaylistId = currentTime
	}

	_, err = s.db.Exec(
		`INSERT INTO playlist (user_id, playlist_id, deleted, title, modified_date, cover_blob)
        VALUES (?, ?, ?, ?, ?, ?)
        ON CONFLICT(user_id, playlist_id) DO UPDATE SET 
            deleted = excluded.deleted,
            title = excluded.title, 
            modified_date = excluded.modified_date, 
            cover_blob = excluded.cover_blob`,
		p.UserId, p.PlaylistId, p.Deleted, p.Title, p.ModifiedDate, p.CoverBlob,
	)

	return p, err
}

func (s *SQLiteStorage) DeletePlaylist(playlistID int64) error {
	user, err := s.GetUser()
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	_, err = s.db.Exec(`DELETE FROM playlist WHERE user_id = ? AND playlist_id = ?`, user.UserId, playlistID)
	return err
}

func (s *SQLiteStorage) GetPlaylist(playlistID int64) (Playlist, error) {
	user, err := s.GetUser()
	if err != nil {
		return Playlist{}, fmt.Errorf("failed to get user: %v", err)
	}

	var p Playlist
	err = s.db.QueryRow(
		`SELECT user_id, playlist_id, deleted, title, modified_date, cover_blob
        FROM playlist WHERE user_id = ? AND playlist_id = ?`,
		user.UserId, playlistID,
	).Scan(&p.UserId, &p.PlaylistId, &p.Deleted, &p.Title, &p.ModifiedDate, &p.CoverBlob)

	return p, err
}

// ---------------- Music Methods ----------------

func (s *SQLiteStorage) PutMusic(m Music) error {
	_, err := s.db.Exec(
		`INSERT INTO music (music_id, source, title, length_seconds) VALUES (?, ?, ?, ?)
        ON CONFLICT(music_id, source) DO UPDATE SET title = excluded.title, length_seconds = excluded.length_seconds`,
		m.MusicId, m.Source, m.Title, m.LengthSeconds,
	)
	return err
}

func (s *SQLiteStorage) DeleteMusic(musicID string, source MusicSource) error {
	_, err := s.db.Exec(`DELETE FROM music WHERE music_id = ? AND source = ?`, musicID, int64(source))
	return err
}

func (s *SQLiteStorage) GetMusic(musicID string, source MusicSource) (Music, error) {
	var m Music
	err := s.db.QueryRow(
		`SELECT music_id, source, title, length_seconds FROM music WHERE music_id = ? AND source = ?`,
		musicID, int64(source),
	).Scan(&m.MusicId, &m.Source, &m.Title, &m.LengthSeconds)

	return m, err
}

// ---------------- File Storer Methods ----------------
// Filesystem storage logic remains unchanged.

func (s *SQLiteStorage) getMusicFilePath(music Music) string {
	return filepath.Join(s.musicFilePath, fmt.Sprintf("%v_%v.mp3", music.Source, music.MusicId))
}

func (s *SQLiteStorage) PutMusicFile(music Music, content io.Reader) error {
	file, err := os.OpenFile(s.getMusicFilePath(music), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0700)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, content)
	file.Close()

	// Remove partially completed files.
	if err != nil {
		os.Remove(s.getMusicFilePath(music))
		return err
	}
	return nil
}

func (s *SQLiteStorage) DeleteMusicFile(music Music) error {
	err := os.Remove(s.getMusicFilePath(music))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (s *SQLiteStorage) GetMusicFile(music Music) (io.ReadCloser, error) {
	return os.Open(s.getMusicFilePath(music))
}

// ---------------- Playlist Manager Methods ----------------

func (s *SQLiteStorage) GetPlaylistsFromUser() ([]Playlist, error) {
	user, err := s.GetUser()
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	rows, err := s.db.Query(
		`SELECT user_id, playlist_id, deleted, title, modified_date, cover_blob 
		 FROM playlist WHERE user_id = ? ORDER BY modified_date DESC`,
		user.UserId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []Playlist
	for rows.Next() {
		var p Playlist
		if err := rows.Scan(&p.UserId, &p.PlaylistId, &p.Deleted, &p.Title, &p.ModifiedDate, &p.CoverBlob); err != nil {
			return nil, err
		}
		playlists = append(playlists, p)
	}
	return playlists, nil
}

func (s *SQLiteStorage) GetMusicFromPlaylist(playlistID int64) ([]Music, error) {
	user, err := s.GetUser()
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	rows, err := s.db.Query(
		`SELECT m.music_id, m.source, m.title, m.length_seconds
        FROM music m
        JOIN playlist_music pm ON m.music_id = pm.music_id AND m.source = pm.source
        WHERE pm.user_id = ? AND pm.playlist_id = ?
        ORDER BY pm.added_at DESC`,
		user.UserId, playlistID,
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

func (s *SQLiteStorage) PutMusicInPlaylist(playlistID int64, musicID string, source MusicSource) error {
	user, err := s.GetUser()
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		`INSERT OR IGNORE INTO playlist_music (user_id, playlist_id, music_id, source, added_at) VALUES (?, ?, ?, ?, ?)`,
		user.UserId, playlistID, musicID, int64(source), time.Now().UnixNano(),
	)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected > 0 {
		if _, err := tx.Exec(`UPDATE playlist SET modified_date = ? WHERE user_id = ? AND playlist_id = ?`, time.Now().UnixNano(), user.UserId, playlistID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SQLiteStorage) DeleteMusicFromPlaylist(playlistID int64, musicID string, source MusicSource) error {
	user, err := s.GetUser()
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		`DELETE FROM playlist_music WHERE user_id = ? AND playlist_id = ? AND music_id = ? AND source = ?`,
		user.UserId, playlistID, musicID, int64(source),
	)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected > 0 {
		if _, err := tx.Exec(`UPDATE playlist SET modified_date = ? WHERE user_id = ? AND playlist_id = ?`, time.Now().UnixNano(), user.UserId, playlistID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
