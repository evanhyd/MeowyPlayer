package storages

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// setupTestDB creates a completely isolated storage instance for each test suite.
func setupTestDB(t *testing.T) *SQLiteStorage {
	t.Helper()

	// Pure in-memory database, NO shared cache.
	// Foreign keys pragma is injected directly into the URI.
	dbPath := "file::memory:?_pragma=foreign_keys(1)"

	// t.TempDir() creates a unique, self-cleaning directory for this specific test
	musicPath := filepath.Join(t.TempDir(), "music")

	storage := NewSQLiteStorage(dbPath, musicPath)
	if storage == nil {
		t.Fatalf("Failed to initialize SQLiteStorage")
	}

	// CRITICAL: Force exactly 1 connection.
	// Without a shared cache, a second connection would point to a blank database.
	storage.db.SetMaxOpenConns(1)
	storage.db.SetMaxIdleConns(1)
	storage.db.SetConnMaxLifetime(0)

	return storage
}

func TestUserOperations(t *testing.T) {
	s := setupTestDB(t)
	defer s.Close()

	profile := UserProfile{
		UserId:           "user_abc123",
		Username:         "test_user",
		Language:         1,
		RegistrationDate: time.Now().Unix(),
		Token:            "token_xyz",
	}

	// 1. Create User
	if err := s.PutUser(profile); err != nil {
		t.Fatalf("PutUser failed: %v", err)
	}

	// 2. Fetch User (Hits cache/DB)
	fetched, err := s.GetUser()
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if fetched.UserId != profile.UserId || fetched.Token != profile.Token {
		t.Errorf("GetUser mismatch. Expected UserId %q, got %q", profile.UserId, fetched.UserId)
	}

	// 3. Update User (Overwrites the single row)
	profile.Username = "updated_name"
	if err := s.PutUser(profile); err != nil {
		t.Fatalf("PutUser (Update) failed: %v", err)
	}
	fetched, _ = s.GetUser()
	if fetched.Username != "updated_name" {
		t.Errorf("Expected username to be updated to 'updated_name', got %q", fetched.Username)
	}

	// 4. Delete User
	if err := s.DeleteUser(); err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}
	if _, err := s.GetUser(); err == nil {
		t.Error("Expected GetUser to fail after DeleteUser, but it succeeded")
	}
}

func TestPlaylistOperations(t *testing.T) {
	s := setupTestDB(t)
	defer s.Close()

	// Seed User Context (Playlists depend on a valid user)
	s.PutUser(UserProfile{UserId: "u1", Username: "user1", Token: "t1"})

	p := Playlist{
		PlaylistId: 1001, // Explicitly set ID
		Title:      "My Favorites",
		Deleted:    false,
		CoverBlob:  []byte{},
	}

	// 1. Create Playlist
	created, err := s.PutPlaylist(p)
	if err != nil {
		t.Fatalf("PutPlaylist failed: %v", err)
	}

	// 2. Fetch Single Playlist
	fetched, err := s.GetPlaylist(created.PlaylistId)
	if err != nil {
		t.Fatalf("GetPlaylist failed: %v", err)
	}
	if fetched.Title != p.Title {
		t.Errorf("Expected title %q, got %q", p.Title, fetched.Title)
	}

	// 3. Fetch All Playlists for User
	// Explicitly set a DIFFERENT ID so it doesn't overwrite the first one
	s.PutPlaylist(Playlist{PlaylistId: 1002, Title: "Workout Mix", CoverBlob: []byte{}})

	playlists, err := s.GetPlaylistsFromUser()
	if err != nil {
		t.Fatalf("GetPlaylistsFromUser failed: %v", err)
	}
	if len(playlists) != 2 {
		t.Errorf("Expected 2 playlists, found %d", len(playlists))
	}

	// 4. Delete Playlist
	if err := s.DeletePlaylist(created.PlaylistId); err != nil {
		t.Fatalf("DeletePlaylist failed: %v", err)
	}
	if _, err := s.GetPlaylist(created.PlaylistId); err == nil {
		t.Error("Expected GetPlaylist to fail after deletion")
	}
}

func TestMusicOperations(t *testing.T) {
	s := setupTestDB(t)
	defer s.Close()

	m := Music{
		MusicId:       "vid_123",
		Source:        YouTubeSource,
		Title:         "Epic Song",
		LengthSeconds: 210,
	}

	// 1. Create Music
	if err := s.PutMusic(m); err != nil {
		t.Fatalf("PutMusic failed: %v", err)
	}

	// 2. Fetch Music
	fetched, err := s.GetMusic(m.MusicId, YouTubeSource)
	if err != nil {
		t.Fatalf("GetMusic failed: %v", err)
	}
	if fetched.Title != m.Title {
		t.Errorf("Expected title %q, got %q", m.Title, fetched.Title)
	}

	// 3. Delete Music
	if err := s.DeleteMusic(m.MusicId, YouTubeSource); err != nil {
		t.Fatalf("DeleteMusic failed: %v", err)
	}
	if _, err := s.GetMusic(m.MusicId, YouTubeSource); err == nil {
		t.Error("Expected GetMusic to fail after deletion")
	}
}

func TestPlaylistMusicOperations(t *testing.T) {
	s := setupTestDB(t)
	defer s.Close()

	// 1. Provision Context (User -> Playlist -> Music)
	s.PutUser(UserProfile{UserId: "u1", Token: "t1"})
	playlist, _ := s.PutPlaylist(Playlist{Title: "Chill", CoverBlob: []byte{}})

	m1 := Music{MusicId: "m1", Source: SpotifySource, Title: "Song A"}
	m2 := Music{MusicId: "m2", Source: YouTubeSource, Title: "Song B"}
	s.PutMusic(m1)
	s.PutMusic(m2)

	// 2. Add Music to Playlist
	if err := s.PutMusicInPlaylist(playlist.PlaylistId, m1.MusicId, SpotifySource); err != nil {
		t.Fatalf("PutMusicInPlaylist (m1) failed: %v", err)
	}
	if err := s.PutMusicInPlaylist(playlist.PlaylistId, m2.MusicId, YouTubeSource); err != nil {
		t.Fatalf("PutMusicInPlaylist (m2) failed: %v", err)
	}

	// 3. Retrieve Music from Playlist
	musics, err := s.GetMusicFromPlaylist(playlist.PlaylistId)
	if err != nil {
		t.Fatalf("GetMusicFromPlaylist failed: %v", err)
	}
	if len(musics) != 2 {
		t.Errorf("Expected 2 songs in playlist, found %d", len(musics))
	}

	// 4. Remove Specific Music from Playlist
	if err := s.DeleteMusicFromPlaylist(playlist.PlaylistId, m1.MusicId, SpotifySource); err != nil {
		t.Fatalf("DeleteMusicFromPlaylist failed: %v", err)
	}

	// Verify only m2 is left
	musics, _ = s.GetMusicFromPlaylist(playlist.PlaylistId)
	if len(musics) != 1 || musics[0].MusicId != m2.MusicId {
		t.Errorf("Expected only Song B (m2) to remain, got: %+v", musics)
	}
}

func TestFileOperations(t *testing.T) {
	s := setupTestDB(t)
	defer s.Close()

	m := Music{
		MusicId: "audio_999",
		Source:  YouTubeSource,
	}
	fileContent := "mock_mp3_binary_data"

	// 1. Write File
	if err := s.PutMusicFile(m, strings.NewReader(fileContent)); err != nil {
		t.Fatalf("PutMusicFile failed: %v", err)
	}

	// 2. Read File Back
	reader, err := s.GetMusicFile(m)
	if err != nil {
		t.Fatalf("GetMusicFile failed: %v", err)
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, reader)
	reader.Close() // Explicit close required

	if buf.String() != fileContent {
		t.Errorf("File content mismatch. Expected %q, got %q", fileContent, buf.String())
	}

	// 3. Delete File
	if err := s.DeleteMusicFile(m); err != nil {
		t.Fatalf("DeleteMusicFile failed: %v", err)
	}

	// Verify File is gone
	if _, err := s.GetMusicFile(m); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected os.ErrNotExist after deletion, got: %v", err)
	}
}
