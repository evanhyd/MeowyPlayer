package storages

import (
	_ "embed"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func newTestStorage(user User) *SQLiteStorage {
	storage := NewSQLiteStorage(":memory:", os.TempDir(), user)

	// Register right away.
	if _, err := storage.db.Exec(`INSERT INTO users VALUES (?, ?, ?, ?)`,
		user.UserId, user.Name, user.Salt, user.HashedPassword); err != nil {
		log.Fatalf("failed to create new user: %v", err)
	}
	return storage
}

func TestPlaylistCRUD(t *testing.T) {
	user := User{UserId: 1, Name: "Alice"}
	storage := newTestStorage(user)
	defer storage.Close()

	// Create
	p, err := storage.CreatePlaylist("Favorites", []byte("cover blob"))
	if err != nil {
		t.Fatalf("CreatePlaylist failed: %v", err)
	}

	if p.PlaylistId == 0 {
		t.Fatalf("expected PlaylistID to be non zero")
	}

	// Get
	got, err := storage.GetPlaylist(p.PlaylistId)
	if err != nil {
		t.Fatalf("GetPlaylist failed: %v", err)
	}
	if got.Title != p.Title {
		t.Errorf("expected title %q, got %q", p.Title, got.Title)
	}

	// Update
	p.Title = "Updated Favorites"
	if err := storage.UpdatePlaylist(p); err != nil {
		t.Fatalf("UpdatePlaylist failed: %v", err)
	}

	got, _ = storage.GetPlaylist(p.PlaylistId)
	if got.Title != p.Title {
		t.Errorf("expected updated title %q, got %q", p.Title, got.Title)
	}

	// GetAll
	list, err := storage.GetAllSortedPlaylists()
	if err != nil {
		t.Fatalf("GetAllPlaylists failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 playlist, got %d", len(list))
	}

	// Delete
	if err := storage.DeletePlaylist(p.PlaylistId); err != nil {
		t.Fatalf("DeletePlaylist failed: %v", err)
	}
	_, err = storage.GetPlaylist(p.PlaylistId)
	if err == nil {
		t.Errorf("expected error after deleting playlist")
	}
}

func TestMusicCRUD(t *testing.T) {
	user := User{UserId: 1, Name: "Alice"}
	storage := newTestStorage(user)
	defer storage.Close()

	m := Music{
		MusicId:       "m1",
		Source:        YouTubeSource,
		Title:         "Song One",
		LengthSeconds: 300,
	}

	if err := storage.CreateMusic(m); err != nil {
		t.Fatalf("CreateMusic failed: %v", err)
	}

	got, err := storage.GetMusic(m.MusicId, m.Source)
	if err != nil {
		t.Fatalf("GetMusic failed: %v", err)
	}
	if got.Title != m.Title {
		t.Errorf("expected title %q, got %q", m.Title, got.Title)
	}

	// Update
	m.Title = "Song One Updated"
	if err := storage.UpdateMusic(m); err != nil {
		t.Fatalf("UpdateMusic failed: %v", err)
	}
	got, _ = storage.GetMusic(m.MusicId, m.Source)
	if got.Title != m.Title {
		t.Errorf("expected updated title %q, got %q", m.Title, got.Title)
	}

	// GetAll
	list, err := storage.GetAllMusic()
	if err != nil {
		t.Fatalf("GetAllMusic failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 music, got %d", len(list))
	}

	// Delete
	if err := storage.DeleteMusic(m.MusicId, m.Source); err != nil {
		t.Fatalf("DeleteMusic failed: %v", err)
	}
	_, err = storage.GetMusic(m.MusicId, m.Source)
	if err == nil {
		t.Errorf("expected error after deleting music")
	}
}

func TestPlaylistMusic(t *testing.T) {
	storage := newTestStorage(User{UserId: 1, Name: "Alice"})
	defer storage.Close()

	// Prepare playlist
	p, err := storage.CreatePlaylist("Favorites", []byte("cover blob"))
	if err != nil {
		t.Fatalf("CreatePlaylist failed: %v", err)
	}

	// Prepare music
	m := Music{MusicId: "m1", Source: YouTubeSource, Title: "Song One", LengthSeconds: 300}
	if err := storage.CreateMusic(m); err != nil {
		t.Fatalf("CreateMusic failed: %v", err)
	}

	// Add music
	if err := storage.AddMusicToPlaylist(p.PlaylistId, m.MusicId, m.Source); err != nil {
		t.Fatalf("AddMusicToPlaylist failed: %v", err)
	}

	// List music -> GetAll
	list, err := storage.GetAllSortedMusicFromPlaylist(p.PlaylistId)
	if err != nil {
		t.Fatalf("GetAllMusicFromPlaylist failed: %v", err)
	}
	if len(list) != 1 || list[0].MusicId != m.MusicId {
		t.Errorf("expected 1 music with ID %q, got %v", m.MusicId, list)
	}

	// Remove music
	if err := storage.RemoveMusicFromPlaylist(p.PlaylistId, m.MusicId, m.Source); err != nil {
		t.Fatalf("RemoveMusicFromPlaylist failed: %v", err)
	}
	list, _ = storage.GetAllSortedMusicFromPlaylist(p.PlaylistId)
	if len(list) != 0 {
		t.Errorf("expected 0 music after removal, got %d", len(list))
	}

	// Edge case: remove from invalid playlist
	if err := storage.RemoveMusicFromPlaylist(9999, m.MusicId, m.Source); err != nil {
		t.Errorf("expected no error removing  music from invalid playlist, got %v", err)
	}

	// Edge case: add music to invalid playlist
	if err := storage.AddMusicToPlaylist(9999, m.MusicId, m.Source); err == nil {
		t.Errorf("expected error adding music to invalid playlist")
	}

	// Edge case: list music in invalid playlist -> GetAll
	music, err := storage.GetAllSortedMusicFromPlaylist(9999)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(music) != 0 {
		t.Errorf("expected empty listing music, got %v", len(music))
	}
}

func TestDeletePlaylistCascadesPlaylistMusic(t *testing.T) {
	storage := newTestStorage(User{UserId: 1, Name: "Alice"})
	defer storage.Close()

	// Create a playlist
	p, err := storage.CreatePlaylist("Cascading Test", []byte("cover blob"))
	if err != nil {
		t.Fatalf("CreatePlaylist failed: %v", err)
	}

	// Create two music tracks
	m1 := Music{MusicId: "m1", Source: YouTubeSource, Title: "Song One", LengthSeconds: 120}
	m2 := Music{MusicId: "m2", Source: YouTubeSource, Title: "Song Two", LengthSeconds: 180}
	if err := storage.CreateMusic(m1); err != nil {
		t.Fatalf("CreateMusic m1 failed: %v", err)
	}
	if err := storage.CreateMusic(m2); err != nil {
		t.Fatalf("CreateMusic m2 failed: %v", err)
	}

	// Add both music to playlist
	if err := storage.AddMusicToPlaylist(p.PlaylistId, m1.MusicId, m1.Source); err != nil {
		t.Fatalf("AddMusicToPlaylist m1 failed: %v", err)
	}
	if err := storage.AddMusicToPlaylist(p.PlaylistId, m2.MusicId, m2.Source); err != nil {
		t.Fatalf("AddMusicToPlaylist m2 failed: %v", err)
	}

	// Verify they exist -> GetAll
	list, err := storage.GetAllSortedMusicFromPlaylist(p.PlaylistId)
	if err != nil {
		t.Fatalf("GetAllMusicFromPlaylist failed: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 music in playlist before delete, got %d", len(list))
	}

	// Delete the playlist
	if err := storage.DeletePlaylist(p.PlaylistId); err != nil {
		t.Fatalf("DeletePlaylist failed: %v", err)
	}

	// Verify the playlist is gone
	if _, err := storage.GetPlaylist(p.PlaylistId); err == nil {
		t.Errorf("expected GetPlaylist to fail after deletion")
	}

	// Check that playlist_music entries were also removed
	rows, err := storage.db.Query(
		`SELECT COUNT(*) FROM playlist_music WHERE user_id = ? AND playlist_id = ?`,
		storage.user.UserId, p.PlaylistId,
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
	storage := newTestStorage(User{UserId: 1, Name: "Alice"})
	defer storage.Close()

	music := Music{Source: YouTubeSource, MusicId: "1234abcd"}
	content := "hello"

	_, err := storage.GetMusicFile(music)
	if err == nil {
		t.Errorf("expected error when get music file")
	}

	err = storage.CreateOrUpdateMusicFile(music, strings.NewReader(content))
	if err != nil {
		t.Fatalf("CreateOrUpdateMusicFile failed: %v", err)
	}

	file, err := storage.GetMusicFile(music)
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

	err = storage.RemoveMusicFile(music)
	if err != nil {
		t.Fatalf("RemoveMusicFile failed: %v", err)
	}
}
