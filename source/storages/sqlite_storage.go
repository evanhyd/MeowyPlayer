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
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		slog.Error("failed to open SQLite database", "error", err)
		return nil
	}
	if _, err := db.Exec(schemaSQL); err != nil {
		slog.Error("failed to create schema", "error", err)
		return nil
	}

	if err := os.MkdirAll(musicFilePath, 0700); err != nil {
		slog.Error("failed to create music file directory", "error", err)
		return nil
	}

	return &SQLiteStorage{
		musicFilePath: musicFilePath,
		db:            db,
	}
}

// ---------------- User Storer Methods ----------------

func (s *SQLiteStorage) PutUser(profile UserProfile) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM user_profile`); err != nil {
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO user_profile (user_id, username, language, registration_date, token) VALUES (?, ?, ?, ?, ?)`,
		profile.UserId, profile.Username, profile.Language, profile.RegistrationDate, profile.Token,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	s.userMutex.Lock()
	s.cachedProfile = nil
	s.userMutex.Unlock()

	return nil
}

func (s *SQLiteStorage) GetUser() (UserProfile, error) {
	s.userMutex.Lock()
	defer s.userMutex.Unlock()

	if s.cachedProfile != nil {
		return *s.cachedProfile, nil
	}

	var p UserProfile
	err := s.db.QueryRow(`SELECT user_id, username, language, registration_date, token FROM user_profile LIMIT 1`).
		Scan(&p.UserId, &p.Username, &p.Language, &p.RegistrationDate, &p.Token)
	if err != nil {
		return UserProfile{}, err
	}

	s.cachedProfile = &p
	return p, nil
}

func (s *SQLiteStorage) DeleteUser() error {
	if _, err := s.db.Exec(`DELETE FROM user_profile`); err != nil {
		return err
	}

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

	p.UserId = user.UserId
	now := time.Now().UnixNano()
	if p.PlaylistId == 0 {
		p.PlaylistId = now
	}
	if p.ModifiedDate == 0 {
		p.ModifiedDate = now
	}

	_, err = s.db.Exec(
		`INSERT INTO playlist (user_id, playlist_id, title, modified_date, cover_blob)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(user_id, playlist_id) DO UPDATE SET 
			title = excluded.title, 
			modified_date = excluded.modified_date, 
			cover_blob = excluded.cover_blob`,
		p.UserId, p.PlaylistId, p.Title, p.ModifiedDate, p.CoverBlob,
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
		`SELECT user_id, playlist_id, title, modified_date, cover_blob
		FROM playlist WHERE user_id = ? AND playlist_id = ?`,
		user.UserId, playlistID,
	).Scan(&p.UserId, &p.PlaylistId, &p.Title, &p.ModifiedDate, &p.CoverBlob)

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

func (s *SQLiteStorage) getMusicFilePath(music Music) string {
	return filepath.Join(s.musicFilePath, fmt.Sprintf("%v_%v.mp3", music.Source, music.MusicId))
}

func (s *SQLiteStorage) PutMusicFile(music Music, content io.Reader) error {
	path := s.getMusicFilePath(music)

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, copyErr := io.Copy(file, content)
	closeErr := file.Close()

	if copyErr != nil || closeErr != nil {
		os.Remove(path) // Cleanup partial/failed file
		if copyErr != nil {
			return copyErr
		}
		return closeErr
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

func (s *SQLiteStorage) GetMusicFile(music Music) (io.ReadSeekCloser, error) {
	return os.Open(s.getMusicFilePath(music))
}

// ---------------- Playlist Manager Methods ----------------

func (s *SQLiteStorage) GetPlaylistsFromUser() ([]Playlist, error) {
	user, err := s.GetUser()
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	rows, err := s.db.Query(
		`SELECT user_id, playlist_id,  title, modified_date, cover_blob 
		 FROM playlist WHERE user_id = ? ORDER BY modified_date DESC`,
		user.UserId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	playlists := make([]Playlist, 0)
	for rows.Next() {
		var p Playlist
		if err := rows.Scan(&p.UserId, &p.PlaylistId, &p.Title, &p.ModifiedDate, &p.CoverBlob); err != nil {
			return nil, err
		}
		playlists = append(playlists, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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

	musics := make([]Music, 0)
	for rows.Next() {
		var m Music
		if err := rows.Scan(&m.MusicId, &m.Source, &m.Title, &m.LengthSeconds); err != nil {
			return nil, err
		}
		musics = append(musics, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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

	now := time.Now().UnixNano()
	res, err := tx.Exec(
		`INSERT OR IGNORE INTO playlist_music (user_id, playlist_id, music_id, source, added_at) VALUES (?, ?, ?, ?, ?)`,
		user.UserId, playlistID, musicID, int64(source), now,
	)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected > 0 {
		if _, err := tx.Exec(`UPDATE playlist SET modified_date = ? WHERE user_id = ? AND playlist_id = ?`, now, user.UserId, playlistID); err != nil {
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
