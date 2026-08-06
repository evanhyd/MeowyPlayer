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

var _ Storage = &SQLiteStorage{}

type SQLiteStorage struct {
	musicFilePath string
	db            *sql.DB
	filesystemMux sync.RWMutex
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

// ---------------- Playlist Methods ----------------

func (s *SQLiteStorage) PutPlaylist(session UserSession, p Playlist) (Playlist, error) {
	currentTime := time.Now().UnixNano()
	p.ModifiedDate = currentTime
	p.UserId = session.UserId

	if p.PlaylistId == 0 {
		p.PlaylistId = currentTime
	}

	_, err := s.db.Exec(
		`INSERT INTO playlists (user_id, playlist_id, title, modified_date, cover_blob)
        VALUES (?, ?, ?, ?, ?)
        ON CONFLICT(user_id, playlist_id) DO UPDATE SET 
            title = excluded.title, 
            modified_date = excluded.modified_date, 
            cover_blob = excluded.cover_blob`,
		p.UserId, p.PlaylistId, p.Title, p.ModifiedDate, p.CoverBlob,
	)

	return p, err
}

func (s *SQLiteStorage) DeletePlaylist(session UserSession, playlistID int64) error {
	_, err := s.db.Exec(`DELETE FROM playlists WHERE user_id = ? AND playlist_id = ?`, session.UserId, playlistID)
	return err
}

func (s *SQLiteStorage) GetPlaylist(session UserSession, playlistID int64) (Playlist, error) {
	p := Playlist{UserId: session.UserId}
	err := s.db.QueryRow(
		`SELECT playlist_id, title, modified_date, cover_blob
        FROM playlists WHERE user_id = ? AND playlist_id = ?`,
		p.UserId, playlistID,
	).Scan(&p.PlaylistId, &p.Title, &p.ModifiedDate, &p.CoverBlob)

	return p, err
}

// ---------------- Music Methods ----------------

func (s *SQLiteStorage) PutMusic(session UserSession, m Music) error {
	_, err := s.db.Exec(
		`INSERT INTO music (music_id, source, title, length_seconds) VALUES (?, ?, ?, ?)
        ON CONFLICT(music_id, source) DO UPDATE SET title = excluded.title, length_seconds = excluded.length_seconds`,
		m.MusicId, m.Source, m.Title, m.LengthSeconds,
	)
	return err
}

func (s *SQLiteStorage) DeleteMusic(session UserSession, musicID string, source MusicSource) error {
	_, err := s.db.Exec(`DELETE FROM music WHERE music_id = ? AND source = ?`, musicID, source)
	return err
}

func (s *SQLiteStorage) GetMusic(session UserSession, musicID string, source MusicSource) (Music, error) {
	var m Music
	err := s.db.QueryRow(
		`SELECT music_id, source, title, length_seconds FROM music WHERE music_id = ? AND source = ?`,
		musicID, source,
	).Scan(&m.MusicId, &m.Source, &m.Title, &m.LengthSeconds)

	return m, err
}

// ---------------- File Storer Methods ----------------

func (s *SQLiteStorage) getMusicFilePath(music Music) string {
	return filepath.Join(s.musicFilePath, fmt.Sprintf("%v_%v.mp3", music.Source, music.MusicId))
}

func (s *SQLiteStorage) PutMusicFile(session UserSession, music Music, content io.Reader) error {
	s.filesystemMux.Lock()
	defer s.filesystemMux.Unlock()

	file, err := os.OpenFile(s.getMusicFilePath(music), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0700)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil
		}
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	return err
}

func (s *SQLiteStorage) DeleteMusicFile(session UserSession, music Music) error {
	s.filesystemMux.Lock()
	defer s.filesystemMux.Unlock()

	err := os.Remove(s.getMusicFilePath(music))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (s *SQLiteStorage) GetMusicFile(session UserSession, music Music) (io.ReadCloser, error) {
	s.filesystemMux.RLock()
	defer s.filesystemMux.RUnlock()
	return os.Open(s.getMusicFilePath(music))
}

// ---------------- Playlist Manager Methods ----------------

func (s *SQLiteStorage) GetPlaylistsFromUser(session UserSession) ([]Playlist, error) {
	rows, err := s.db.Query(
		`SELECT playlist_id, title, modified_date, cover_blob FROM playlists WHERE user_id = ? ORDER BY modified_date DESC`,
		session.UserId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []Playlist
	for rows.Next() {
		p := Playlist{UserId: session.UserId}
		if err := rows.Scan(&p.PlaylistId, &p.Title, &p.ModifiedDate, &p.CoverBlob); err != nil {
			return nil, err
		}
		playlists = append(playlists, p)
	}
	return playlists, nil
}

func (s *SQLiteStorage) GetMusicFromPlaylist(session UserSession, playlistID int64) ([]Music, error) {
	rows, err := s.db.Query(
		`SELECT m.music_id, m.source, m.title, m.length_seconds
        FROM music m
        JOIN playlist_music pm ON m.music_id = pm.music_id AND m.source = pm.source
        WHERE pm.user_id = ? AND pm.playlist_id = ?
        ORDER BY pm.added_at DESC`,
		session.UserId, playlistID,
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

func (s *SQLiteStorage) PutMusicInPlaylist(session UserSession, playlistID int64, musicID string, source MusicSource) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		`INSERT OR IGNORE INTO playlist_music (user_id, playlist_id, music_id, source, added_at) VALUES (?, ?, ?, ?, ?)`,
		session.UserId, playlistID, musicID, source, time.Now().UnixNano(),
	)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected > 0 {
		if _, err := tx.Exec(`UPDATE playlists SET modified_date = ? WHERE user_id = ? AND playlist_id = ?`, time.Now().UnixNano(), session.UserId, playlistID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SQLiteStorage) DeleteMusicFromPlaylist(session UserSession, playlistID int64, musicID string, source MusicSource) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		`DELETE FROM playlist_music WHERE user_id = ? AND playlist_id = ? AND music_id = ? AND source = ?`,
		session.UserId, playlistID, musicID, source,
	)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected > 0 {
		if _, err := tx.Exec(`UPDATE playlists SET modified_date = ? WHERE user_id = ? AND playlist_id = ?`, time.Now().UnixNano(), session.UserId, playlistID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
