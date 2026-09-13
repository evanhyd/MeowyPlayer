package storages

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) (*SQLiteStorage, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "sqlite_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tmpDir, "test.db")
	musicPath := filepath.Join(tmpDir, "music")

	storage := NewSQLiteStorage(dbPath, musicPath)
	if storage == nil {
		t.Fatalf("failed to create SQLiteStorage")
	}

	return storage, func() {
		storage.Close()
		os.RemoveAll(tmpDir)
	}
}

func TestSQLiteStorage_User(t *testing.T) {
	s, cleanup := setupTestDB(t)
	defer cleanup()

	profile := UserProfile{
		UserId:           "user_123",
		Username:         "TestUser",
		Language:         LangEnglish,
		RegistrationDate: 1600000000,
		Token:            "token_abc",
	}

	if err := s.PutUser(profile); err != nil {
		t.Fatalf("PutUser failed: %v", err)
	}

	got, err := s.GetUser()
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if !reflect.DeepEqual(profile, got) {
		t.Errorf("GetUser mismatch. got %v, want %v", got, profile)
	}

	if err := s.DeleteUser(); err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}

	if _, err := s.GetUser(); err == nil {
		t.Error("GetUser expected error after deletion, got nil")
	}
}

func TestSQLiteStorage_Playlist(t *testing.T) {
	s, cleanup := setupTestDB(t)
	defer cleanup()

	// Seed user required for foreign keys & GetUser() lookups
	_ = s.PutUser(UserProfile{UserId: "user_123"})

	p := Playlist{
		Title:        "My Favorites",
		CoverBlob:    []byte(""),
		ModifiedDate: time.Now().UnixNano(),
	}

	// Test PutPlaylist (auto-generates IDs and Dates)
	saved, err := s.PutPlaylist(p)
	if err != nil {
		t.Fatalf("PutPlaylist failed: %v", err)
	}
	if saved.PlaylistId == 0 {
		t.Error("PutPlaylist failed to auto-generate PlaylistId")
	}
	if saved.UserId != "user_123" {
		t.Errorf("PutPlaylist failed to enforce UserId. got %v", saved.UserId)
	}

	// Test GetPlaylist
	got, err := s.GetPlaylist(saved.PlaylistId)
	if err != nil {
		t.Fatalf("GetPlaylist failed: %v", err)
	}
	if got.Title != p.Title {
		t.Errorf("GetPlaylist mismatch. got title %v, want %v", got.Title, p.Title)
	}

	// Test GetPlaylists
	playlists, err := s.GetPlaylists()
	if err != nil {
		t.Fatalf("GetPlaylists failed: %v", err)
	}
	if len(playlists) != 1 {
		t.Errorf("GetPlaylists expected 1, got %d", len(playlists))
	}

	// Test DeletePlaylist
	if err := s.DeletePlaylist(saved.PlaylistId); err != nil {
		t.Fatalf("DeletePlaylist failed: %v", err)
	}
	if _, err := s.GetPlaylist(saved.PlaylistId); err == nil {
		t.Error("GetPlaylist expected error after deletion, got nil")
	}
}

func TestSQLiteStorage_Music(t *testing.T) {
	s, cleanup := setupTestDB(t)
	defer cleanup()

	m := Music{
		MusicId:       "vid_123",
		Source:        YouTubeSource,
		Title:         "Epic Song",
		LengthSeconds: 210,
	}

	if err := s.PutMusic(m); err != nil {
		t.Fatalf("PutMusic failed: %v", err)
	}

	got, err := s.GetMusic(m.MusicId, m.Source)
	if err != nil {
		t.Fatalf("GetMusic failed: %v", err)
	}
	if !reflect.DeepEqual(m, got) {
		t.Errorf("GetMusic mismatch. got %v, want %v", got, m)
	}

	// Test Conflict Update (UPSERT)
	m.Title = "Epic Song (Remastered)"
	if err := s.PutMusic(m); err != nil {
		t.Fatalf("PutMusic UPSERT failed: %v", err)
	}

	gotUpdate, _ := s.GetMusic(m.MusicId, m.Source)
	if gotUpdate.Title != "Epic Song (Remastered)" {
		t.Errorf("PutMusic UPSERT mismatch. got %v", gotUpdate.Title)
	}

	if err := s.DeleteMusic(m.MusicId, m.Source); err != nil {
		t.Fatalf("DeleteMusic failed: %v", err)
	}
}

func TestSQLiteStorage_PlaylistMusic(t *testing.T) {
	s, cleanup := setupTestDB(t)
	defer cleanup()

	_ = s.PutUser(UserProfile{UserId: "user_123"})
	pl, err := s.PutPlaylist(Playlist{Title: "Mix", CoverBlob: []byte("")})
	if err != nil {
		t.Fatalf(err.Error())
	}
	m1 := Music{MusicId: "m1", Source: YouTubeSource}
	m2 := Music{MusicId: "m2", Source: SpotifySource}
	_ = s.PutMusic(m1)
	_ = s.PutMusic(m2)

	rel1 := PlaylistMusic{
		PlaylistId:   pl.PlaylistId,
		MusicId:      m1.MusicId,
		Source:       int64(m1.Source),
		ModifiedDate: 100,
	}
	rel2 := PlaylistMusic{
		PlaylistId:   pl.PlaylistId,
		MusicId:      m2.MusicId,
		Source:       int64(m2.Source),
		ModifiedDate: 200,
	}

	if err := s.PutMusic(m1); err != nil {
		t.Fatalf("PutMusic failed: %v", err)
	}

	if err := s.PutMusic(m2); err != nil {
		t.Fatalf("PutMusic failed: %v", err)
	}

	if err := s.PutPlaylistMusic(rel1); err != nil {
		t.Fatalf("PutPlaylistMusic failed: %v", err)
	}
	if err := s.PutPlaylistMusic(rel2); err != nil {
		t.Fatalf("PutPlaylistMusic failed: %v", err)
	}

	// Test GetPlaylistMusic order
	rels, err := s.GetAllPlaylistMusic(pl.PlaylistId)
	if err != nil {
		t.Fatalf("GetPlaylistMusic failed: %v", err)
	}
	if len(rels) != 2 {
		t.Fatalf("Expected 2 relations, got %d", len(rels))
	}
	if rels[0].MusicId != "m2" || rels[1].MusicId != "m1" {
		t.Errorf("GetPlaylistMusic returned incorrect order")
	}

	// Test DeletePlaylistMusic
	toDelete := PlaylistMusic{PlaylistId: pl.PlaylistId, MusicId: m1.MusicId, Source: m1.Source, ModifiedDate: time.Now().UnixNano()}
	if err := s.DeletePlaylistMusic(toDelete); err != nil {
		t.Fatalf("DeletePlaylistMusic failed: %v", err)
	}

	relsAfter, _ := s.GetAllPlaylistMusic(pl.PlaylistId)
	if len(relsAfter) != 1 {
		t.Errorf("Expected 1 relation after deletion, got %d", len(relsAfter))
	}
}

func TestSQLiteStorage_File(t *testing.T) {
	s, cleanup := setupTestDB(t)
	defer cleanup()

	m := Music{
		MusicId: "audio_1",
		Source:  YouTubeSource,
	}
	content := []byte("fake_mp3_data")

	// Test Put
	if err := s.PutMusicFile(m, bytes.NewReader(content)); err != nil {
		t.Fatalf("PutMusicFile failed: %v", err)
	}

	// Test Get
	reader, err := s.GetMusicFile(m)
	if err != nil {
		t.Fatalf("GetMusicFile failed: %v", err)
	}

	readData, err := io.ReadAll(reader)

	// EXPLICITLY CLOSE the reader right after we finish reading.
	// Do not use 'defer' here, otherwise Windows will block the os.Remove call below.
	if errClose := reader.Close(); errClose != nil {
		t.Fatalf("Failed to close file reader: %v", errClose)
	}

	if err != nil {
		t.Fatalf("Failed to read file content: %v", err)
	}
	if !bytes.Equal(readData, content) {
		t.Errorf("File content mismatch. got %s, want %s", readData, content)
	}

	// Test Delete - will now succeed because the file handle was closed
	if err := s.DeleteMusicFile(m); err != nil {
		t.Fatalf("DeleteMusicFile failed: %v", err)
	}

	if _, err := s.GetMusicFile(m); err == nil {
		t.Error("GetMusicFile expected error after deletion, got nil")
	}
}
