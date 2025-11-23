package storage

import (
	_ "embed"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func newTestManager(user User) *SQLiteStorageManager {
	storageManager := NewSQLiteStorageManager(":memory:", os.TempDir(), user)

	// Register right away.
	if _, err := storageManager.db.Exec(`INSERT INTO users VALUES (?, ?, ?, ?)`,
		user.UserID, user.Name, user.Salt, user.HashedPassword); err != nil {
		log.Fatalf("failed to create new user: %v", err)
	}
	return storageManager
}

func TestPlaylistCRUD(t *testing.T) {
	user := User{UserID: 1, Name: "Alice"}
	manager := newTestManager(user)
	defer manager.Close()

	p := Playlist{
		Title:     "Favorites",
		CoverBlob: []byte("cover blob"),
	}

	// Create
	p, err := manager.CreatePlaylist(p)
	if err != nil {
		t.Fatalf("CreatePlaylist failed: %v", err)
	}

	if p.PlaylistID == 0 {
		t.Fatalf("expected PlaylistID to be non zero")
	}

	// Get
	got, err := manager.GetPlaylist(p.PlaylistID)
	if err != nil {
		t.Fatalf("GetPlaylist failed: %v", err)
	}
	if got.Title != p.Title {
		t.Errorf("expected title %q, got %q", p.Title, got.Title)
	}

	// Update
	p.Title = "Updated Favorites"
	if err := manager.UpdatePlaylist(p); err != nil {
		t.Fatalf("UpdatePlaylist failed: %v", err)
	}

	got, _ = manager.GetPlaylist(p.PlaylistID)
	if got.Title != p.Title {
		t.Errorf("expected updated title %q, got %q", p.Title, got.Title)
	}

	// List
	list, err := manager.ListAllPlaylists()
	if err != nil {
		t.Fatalf("ListPlaylists failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 playlist, got %d", len(list))
	}

	// Delete
	if err := manager.DeletePlaylist(p.PlaylistID); err != nil {
		t.Fatalf("DeletePlaylist failed: %v", err)
	}
	_, err = manager.GetPlaylist(p.PlaylistID)
	if err == nil {
		t.Errorf("expected error after deleting playlist")
	}
}

func TestMusicCRUD(t *testing.T) {
	user := User{UserID: 1, Name: "Alice"}
	manager := newTestManager(user)
	defer manager.Close()

	m := Music{
		MusicID:       "m1",
		Source:        YouTubeSource,
		Title:         "Song One",
		LengthSeconds: 300,
	}

	if err := manager.CreateMusic(m); err != nil {
		t.Fatalf("CreateMusic failed: %v", err)
	}

	got, err := manager.GetMusic(m.MusicID, m.Source)
	if err != nil {
		t.Fatalf("GetMusic failed: %v", err)
	}
	if got.Title != m.Title {
		t.Errorf("expected title %q, got %q", m.Title, got.Title)
	}

	// Update
	m.Title = "Song One Updated"
	if err := manager.UpdateMusic(m); err != nil {
		t.Fatalf("UpdateMusic failed: %v", err)
	}
	got, _ = manager.GetMusic(m.MusicID, m.Source)
	if got.Title != m.Title {
		t.Errorf("expected updated title %q, got %q", m.Title, got.Title)
	}

	// List
	list, err := manager.ListAllMusic()
	if err != nil {
		t.Fatalf("ListMusic failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 music, got %d", len(list))
	}

	// Delete
	if err := manager.DeleteMusic(m.MusicID, m.Source); err != nil {
		t.Fatalf("DeleteMusic failed: %v", err)
	}
	_, err = manager.GetMusic(m.MusicID, m.Source)
	if err == nil {
		t.Errorf("expected error after deleting music")
	}
}

func TestPlaylistMusic(t *testing.T) {
	manager := newTestManager(User{UserID: 1, Name: "Alice"})
	defer manager.Close()

	// Prepare playlist
	p := Playlist{
		Title:     "Favorites",
		CoverBlob: []byte("cover blob"),
	}

	p, err := manager.CreatePlaylist(p)
	if err != nil {
		t.Fatalf("CreatePlaylist failed: %v", err)
	}

	// Prepare music
	m := Music{MusicID: "m1", Source: YouTubeSource, Title: "Song One", LengthSeconds: 300}
	if err := manager.CreateMusic(m); err != nil {
		t.Fatalf("CreateMusic failed: %v", err)
	}

	// Add music
	if err := manager.AddMusicToPlaylist(p.PlaylistID, m.MusicID, m.Source); err != nil {
		t.Fatalf("AddMusicToPlaylist failed: %v", err)
	}

	// List music
	list, err := manager.ListMusicInPlaylist(p.PlaylistID)
	if err != nil {
		t.Fatalf("ListMusicInPlaylist failed: %v", err)
	}
	if len(list) != 1 || list[0].MusicID != m.MusicID {
		t.Errorf("expected 1 music with ID %q, got %v", m.MusicID, list)
	}

	// Remove music
	if err := manager.RemoveMusicFromPlaylist(p.PlaylistID, m.MusicID, m.Source); err != nil {
		t.Fatalf("RemoveMusicFromPlaylist failed: %v", err)
	}
	list, _ = manager.ListMusicInPlaylist(p.PlaylistID)
	if len(list) != 0 {
		t.Errorf("expected 0 music after removal, got %d", len(list))
	}

	// Edge case: remove from invalid playlist
	if err := manager.RemoveMusicFromPlaylist(9999, m.MusicID, m.Source); err != nil {
		t.Errorf("expected no error removing  music from invalid playlist, got %v", err)
	}

	// Edge case: add music to invalid playlist
	if err := manager.AddMusicToPlaylist(9999, m.MusicID, m.Source); err == nil {
		t.Errorf("expected error adding music to invalid playlist")
	}

	// Edge case: list music in invalid playlist
	music, err := manager.ListMusicInPlaylist(9999)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(music) != 0 {
		t.Errorf("expected empty listing music, got %v", len(music))
	}
}

func TestDeletePlaylistCascadesPlaylistMusic(t *testing.T) {
	manager := newTestManager(User{UserID: 1, Name: "Alice"})
	defer manager.Close()

	// Create a playlist
	p, err := manager.CreatePlaylist(Playlist{Title: "Cascading Test", CoverBlob: []byte("cover blob")})
	if err != nil {
		t.Fatalf("CreatePlaylist failed: %v", err)
	}

	// Create two music tracks
	m1 := Music{MusicID: "m1", Source: YouTubeSource, Title: "Song One", LengthSeconds: 120}
	m2 := Music{MusicID: "m2", Source: YouTubeSource, Title: "Song Two", LengthSeconds: 180}
	if err := manager.CreateMusic(m1); err != nil {
		t.Fatalf("CreateMusic m1 failed: %v", err)
	}
	if err := manager.CreateMusic(m2); err != nil {
		t.Fatalf("CreateMusic m2 failed: %v", err)
	}

	// Add both music to playlist
	if err := manager.AddMusicToPlaylist(p.PlaylistID, m1.MusicID, m1.Source); err != nil {
		t.Fatalf("AddMusicToPlaylist m1 failed: %v", err)
	}
	if err := manager.AddMusicToPlaylist(p.PlaylistID, m2.MusicID, m2.Source); err != nil {
		t.Fatalf("AddMusicToPlaylist m2 failed: %v", err)
	}

	// Verify they exist
	list, err := manager.ListMusicInPlaylist(p.PlaylistID)
	if err != nil {
		t.Fatalf("ListMusicInPlaylist failed: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 music in playlist before delete, got %d", len(list))
	}

	// Delete the playlist
	if err := manager.DeletePlaylist(p.PlaylistID); err != nil {
		t.Fatalf("DeletePlaylist failed: %v", err)
	}

	// Verify the playlist is gone
	if _, err := manager.GetPlaylist(p.PlaylistID); err == nil {
		t.Errorf("expected GetPlaylist to fail after deletion")
	}

	// Check that playlist_music entries were also removed
	rows, err := manager.db.Query(
		`SELECT COUNT(*) FROM playlist_music WHERE user_id = ? AND playlist_id = ?`,
		manager.user.UserID, p.PlaylistID,
	)
	if err != nil {
		t.Fatalf("query playlist_music count failed: %v", err)
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			t.Fatalf("scan count failed: %v", err)
		}
	}
	if count != 0 {
		t.Errorf("expected 0 playlist_music after playlist delete, got %d", count)
	}
}

func TestReadWriteMusicFile(t *testing.T) {
	manager := newTestManager(User{UserID: 1, Name: "Alice"})
	defer manager.Close()

	music := Music{Source: YouTubeSource, MusicID: "1234abcd"}
	content := "hello"

	_, err := manager.GetMusicFile(music)
	if err == nil {
		t.Errorf("expected error when get music file")
	}

	err = manager.CreateOrUpdateMusicFile(music, strings.NewReader(content))
	if err != nil {
		t.Fatalf("CreateOrUpdateMusicFile failed: %v", err)
	}

	file, err := manager.GetMusicFile(music)
	if err != nil {
		t.Fatalf("GetMusicFile failed: %v", err)
	}

	readContent, err := io.ReadAll(file)
	if err != nil {
		file.Close()
		t.Fatalf("Read music file content failed: %v", err)
	}
	if string(readContent) != content {
		file.Close()
		t.Fatalf("Expected content = %v, got %v", content, string(readContent))
	}
	file.Close()

	err = manager.RemoveMusicFile(music)
	if err != nil {
		t.Fatalf("RemoveMusicFile failed: %v", err)
	}
}
