package storages

import (
	"bytes"
	"io"
	"path/filepath"
	"testing"
	"time"
)

func setupTestStorage(t *testing.T) *SQLiteStorage {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	musicDir := filepath.Join(tempDir, "music")

	storage := NewSQLiteStorage(dbPath, musicDir)
	if storage == nil {
		t.Fatalf("Failed to initialize SQLiteStorage")
	}

	t.Cleanup(func() {
		storage.Close()
	})

	return storage
}

func createTestUser(t *testing.T, s *SQLiteStorage) int64 {
	userID := int64(1)
	_, err := s.db.Exec(`INSERT INTO users (user_id, name, hashed_password, salt) VALUES (?, 'Test User', '', '')`, userID)
	if err != nil {
		t.Fatalf("Failed to insert dummy test user: %v", err)
	}
	return userID
}

func getTestUserSession(userID int64) UserSession {
	return UserSession{UserId: userID}
}

func TestMusicStorer_Idempotency(t *testing.T) {
	s := setupTestStorage(t)
	userID := createTestUser(t, s)
	session := getTestUserSession(userID)

	m := Music{
		MusicId:       "m1",
		Source:        YouTubeSource,
		Title:         "Original Title",
		LengthSeconds: 120,
	}

	// 1. Initial Put
	if err := s.PutMusic(session, m); err != nil {
		t.Fatalf("Failed initial PutMusic: %v", err)
	}

	// 2. Idempotent Put (Update)
	m.Title = "Updated Title"
	if err := s.PutMusic(session, m); err != nil {
		t.Fatalf("Failed idempotent PutMusic: %v", err)
	}

	// 3. Verify state
	fetched, err := s.GetMusic(session, m.MusicId, m.Source)
	if err != nil {
		t.Fatalf("Failed GetMusic: %v", err)
	}
	if fetched.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", fetched.Title)
	}

	// 4. Initial Delete
	if err := s.DeleteMusic(session, m.MusicId, m.Source); err != nil {
		t.Fatalf("Failed initial DeleteMusic: %v", err)
	}

	// 5. Idempotent Delete (Should not return error)
	if err := s.DeleteMusic(session, m.MusicId, m.Source); err != nil {
		t.Fatalf("Failed idempotent DeleteMusic: %v", err)
	}
}

func TestPlaylistStorer_Idempotency(t *testing.T) {
	s := setupTestStorage(t)
	userID := createTestUser(t, s)
	session := getTestUserSession(userID)

	p := Playlist{
		Title:     "My First Playlist",
		CoverBlob: []byte("fake-image-data"),
	}

	// 1. Initial Put (ID generation)
	createdP, err := s.PutPlaylist(session, p)
	if err != nil {
		t.Fatalf("Failed initial PutPlaylist: %v", err)
	}
	if createdP.PlaylistId == 0 {
		t.Fatal("Expected PlaylistId to be generated, got 0")
	}

	// 2. Idempotent Put (Update)
	createdP.Title = "My Updated Playlist"
	updatedP, err := s.PutPlaylist(session, createdP)
	if err != nil {
		t.Fatalf("Failed idempotent PutPlaylist: %v", err)
	}
	if updatedP.PlaylistId != createdP.PlaylistId {
		t.Fatalf("Playlist ID changed during update")
	}

	// 3. Verify State
	fetched, err := s.GetPlaylist(session, updatedP.PlaylistId)
	if err != nil {
		t.Fatalf("Failed GetPlaylist: %v", err)
	}
	if fetched.Title != "My Updated Playlist" {
		t.Errorf("Expected title 'My Updated Playlist', got '%s'", fetched.Title)
	}

	// 4. Initial Delete
	if err := s.DeletePlaylist(session, fetched.PlaylistId); err != nil {
		t.Fatalf("Failed initial DeletePlaylist: %v", err)
	}

	// 5. Idempotent Delete
	if err := s.DeletePlaylist(session, fetched.PlaylistId); err != nil {
		t.Fatalf("Failed idempotent DeletePlaylist: %v", err)
	}
}

func TestPlaylistLinker_Idempotency(t *testing.T) {
	s := setupTestStorage(t)
	userID := createTestUser(t, s)
	session := getTestUserSession(userID)

	// Setup: Create a playlist and a music track first
	p, err := s.PutPlaylist(session, Playlist{Title: "Link Test", CoverBlob: []byte("")})
	if err != nil {
		t.Fatalf("failed to put a playlist: %v", err)
	}
	s.PutMusic(session, Music{MusicId: "m1", Source: SpotifySource, Title: "Track 1"})

	// 1. Initial Link
	if err := s.PutMusicInPlaylist(session, p.PlaylistId, "m1", SpotifySource); err != nil {
		t.Fatalf("Failed initial PutMusicInPlaylist: %v", err)
	}

	// 2. Idempotent Link (Should not violate unique constraint)
	if err := s.PutMusicInPlaylist(session, p.PlaylistId, "m1", SpotifySource); err != nil {
		t.Fatalf("Failed idempotent PutMusicInPlaylist: %v", err)
	}

	// 3. Verify Links
	tracks, err := s.GetMusicFromPlaylist(session, p.PlaylistId)
	if err != nil {
		t.Fatalf("Failed GetAllSortedMusicFromPlaylist: %v", err)
	}
	if len(tracks) != 1 {
		t.Fatalf("Expected 1 track in playlist, got %d", len(tracks))
	}

	// 4. Initial Unlink
	if err := s.DeleteMusicFromPlaylist(session, p.PlaylistId, "m1", SpotifySource); err != nil {
		t.Fatalf("Failed initial DeleteMusicFromPlaylist: %v", err)
	}

	// 5. Idempotent Unlink
	if err := s.DeleteMusicFromPlaylist(session, p.PlaylistId, "m1", SpotifySource); err != nil {
		t.Fatalf("Failed idempotent DeleteMusicFromPlaylist: %v", err)
	}
}

func TestFileStorer_Idempotency(t *testing.T) {
	s := setupTestStorage(t)
	userID := createTestUser(t, s)
	session := getTestUserSession(userID)

	m := Music{MusicId: "file1", Source: YouTubeSource}
	content1 := []byte("audio-data-v1")
	content2 := []byte("audio-data-v2")

	// 1. Initial File Put
	if err := s.PutMusicFile(session, m, bytes.NewReader(content1)); err != nil {
		t.Fatalf("Failed initial PutMusicFile: %v", err)
	}

	// 2. Idempotent File Put (Overwrite)
	if err := s.PutMusicFile(session, m, bytes.NewReader(content2)); err != nil {
		t.Fatalf("Failed idempotent PutMusicFile: %v", err)
	}

	// 3. Read and Verify (Should match content2)
	rc, err := s.GetMusicFile(session, m)
	if err != nil {
		t.Fatalf("Failed GetMusicFile: %v", err)
	}

	readData, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("Failed to read data: %v", err)
	}
	if string(readData) != string(content2) {
		t.Errorf("Expected file content %s, got %s", content2, readData)
		rc.Close()
	}
	rc.Close()

	// 4. Initial Delete
	if err := s.DeleteMusicFile(session, m); err != nil {
		t.Fatalf("Failed initial DeleteMusicFile: %v", err)
	}

	// 5. Idempotent Delete (File already gone, should not error)
	if err := s.DeleteMusicFile(session, m); err != nil {
		t.Fatalf("Failed idempotent DeleteMusicFile: %v", err)
	}
}

func TestGetPlaylistsFromUser(t *testing.T) {
	s := setupTestStorage(t)
	userID := createTestUser(t, s)
	session := getTestUserSession(userID)

	// 1. Test Empty State
	playlists, err := s.GetPlaylistsFromUser(session)
	if err != nil {
		t.Fatalf("Expected no error for empty state, got %v", err)
	}
	if len(playlists) != 0 {
		t.Fatalf("Expected 0 playlists, got %d", len(playlists))
	}

	// 2. Insert Test Data
	p1, err := s.PutPlaylist(session, Playlist{
		Title:     "Older Playlist",
		CoverBlob: []byte{},
	})
	if err != nil {
		t.Fatalf("Failed to insert first playlist: %v", err)
	}

	// Guarantee a distinct timestamp for the descending order check
	time.Sleep(1 * time.Millisecond)

	p2, err := s.PutPlaylist(session, Playlist{
		Title:     "Newer Playlist",
		CoverBlob: []byte{},
	})
	if err != nil {
		t.Fatalf("Failed to insert second playlist: %v", err)
	}

	// 3. Test Retrieval and Sorting (ORDER BY modified_date DESC)
	playlists, err = s.GetPlaylistsFromUser(session)
	if err != nil {
		t.Fatalf("Failed to GetPlaylistsFromUser: %v", err)
	}
	if len(playlists) != 2 {
		t.Fatalf("Expected 2 playlists, got %d", len(playlists))
	}

	// p2 should be first because it was created last
	if playlists[0].PlaylistId != p2.PlaylistId || playlists[0].Title != "Newer Playlist" {
		t.Errorf("Expected first playlist to be 'Newer Playlist', got '%s'", playlists[0].Title)
	}
	if playlists[1].PlaylistId != p1.PlaylistId || playlists[1].Title != "Older Playlist" {
		t.Errorf("Expected second playlist to be 'Older Playlist', got '%s'", playlists[1].Title)
	}
}
