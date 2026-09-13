package storages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"meowyplayer/schemas"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

type mockEndpointProvider map[string]string

func (m mockEndpointProvider) GetEndpoint(key string) (string, error) {
	if v, ok := m[key]; ok && v != "" {
		return v, nil
	}
	return "", fmt.Errorf("endpoint %v is not configured", key)
}

func setupTestServerStorage(t *testing.T) (*ServerStorage, *httptest.Server, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "server_storage_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tmpDir, "test.db")
	musicPath := filepath.Join(tmpDir, "music")
	localDB := NewSQLiteStorage(dbPath, musicPath)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{}`)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))

	mockEndpoints := mockEndpointProvider{
		"me":                      ts.URL + "/me",
		"putPlaylist":             ts.URL + "/putPlaylist",
		"putMusicBulk":            ts.URL + "/putMusicBulk",
		"putMusicInPlaylistBulk":  ts.URL + "/putMusicInPlaylistBulk",
		"getPlaylistContent":      ts.URL + "/getPlaylistContent",
		"getPlaylistsFromUser":    ts.URL + "/getPlaylistsFromUser",
		"deletePlaylist":          ts.URL + "/deletePlaylist",
		"putMusic":                ts.URL + "/putMusic",
		"putMusicInPlaylist":      ts.URL + "/putMusicInPlaylist",
		"deleteMusicFromPlaylist": ts.URL + "/deleteMusicFromPlaylist",
	}

	serverStorage := NewServerStorage(ts.Client(), localDB, mockEndpoints)

	return serverStorage, ts, func() {
		if err := serverStorage.Close(); err != nil {
			t.Errorf("failed to close server storage: %v", err)
		}
		ts.Close()
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("failed to remove temporary test directory: %v", err)
		}
	}
}

// ---------------- User Storer Tests ----------------

func TestServerStorage_UserOps(t *testing.T) {
	s, _, cleanup := setupTestServerStorage(t)
	defer cleanup()

	profile := UserProfile{UserId: "user_123", Username: "Test", Token: "token_abc"}

	if err := s.PutUser(profile); err != nil {
		t.Fatalf("PutUser failed: %v", err)
	}

	got, err := s.GetUser()
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if got.UserId != profile.UserId {
		t.Errorf("GetUser mismatch. got %v, want %v", got.UserId, profile.UserId)
	}

	s.setSynced(true)
	if err := s.DeleteUser(); err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}
	if s.hasSynced() {
		t.Error("DeleteUser failed to reset synced flag")
	}
}

// ---------------- Playlist Storer Tests ----------------

func TestServerStorage_PlaylistOps(t *testing.T) {
	s, _, cleanup := setupTestServerStorage(t)
	defer cleanup()

	if err := s.PutUser(UserProfile{UserId: "user_123", Token: "mock_token"}); err != nil {
		t.Fatalf("failed seeding user: %v", err)
	}

	p := Playlist{Title: "My Playlist", CoverBlob: []byte("img")}

	saved, err := s.PutPlaylist(p)
	if err != nil {
		t.Fatalf("PutPlaylist failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	got, err := s.GetPlaylist(saved.PlaylistId)
	if err != nil {
		t.Fatalf("GetPlaylist failed: %v", err)
	}
	if got.Title != p.Title {
		t.Errorf("Title mismatch. got %v, want %v", got.Title, p.Title)
	}

	s.setSynced(true)
	lists, err := s.GetPlaylists()
	if err != nil {
		t.Fatalf("GetPlaylists failed: %v", err)
	}
	if len(lists) != 1 {
		t.Errorf("Expected 1 playlist, got %d", len(lists))
	}

	if err := s.DeletePlaylist(saved.PlaylistId); err != nil {
		t.Fatalf("DeletePlaylist failed: %v", err)
	}
	time.Sleep(50 * time.Millisecond)
}

// ---------------- Music Storer Tests ----------------

func TestServerStorage_MusicOps(t *testing.T) {
	s, _, cleanup := setupTestServerStorage(t)
	defer cleanup()

	m := Music{MusicId: "m1", Source: YouTubeSource, Title: "Song", LengthSeconds: 100}

	if err := s.PutMusic(m); err != nil {
		t.Fatalf("PutMusic failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	got, err := s.GetMusic(m.MusicId, m.Source)
	if err != nil {
		t.Fatalf("GetMusic failed: %v", err)
	}
	if !reflect.DeepEqual(m, got) {
		t.Errorf("GetMusic mismatch. got %v, want %v", got, m)
	}

	if err := s.DeleteMusic(m.MusicId, m.Source); err != nil {
		t.Fatalf("DeleteMusic failed: %v", err)
	}
}

// ---------------- PlaylistMusic Storer Tests ----------------

func TestServerStorage_PlaylistMusicOps(t *testing.T) {
	s, _, cleanup := setupTestServerStorage(t)
	defer cleanup()

	if err := s.PutUser(UserProfile{UserId: "user_123", Token: "tok"}); err != nil {
		t.Fatalf("PutUser err: %v", err)
	}
	p, err := s.Storage.PutPlaylist(Playlist{Title: "List", CoverBlob: []byte{}})
	if err != nil {
		t.Fatalf("PutPlaylist err: %v", err)
	}

	// FIX: Must insert the music into the database first to satisfy the FOREIGN KEY constraint
	if err := s.PutMusic(Music{MusicId: "m1", Source: YouTubeSource, Title: "Song", LengthSeconds: 100}); err != nil {
		t.Fatalf("PutMusic failed: %v", err)
	}

	rel := PlaylistMusic{
		PlaylistId: p.PlaylistId,
		MusicId:    "m1",
		Source:     int64(YouTubeSource),
		AddedAt:    12345,
	}

	if err := s.PutPlaylistMusic(rel); err != nil {
		t.Fatalf("PutPlaylistMusic failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	rels, err := s.GetAllPlaylistMusic(p.PlaylistId)
	if err != nil {
		t.Fatalf("GetPlaylistMusic failed: %v", err)
	}
	if len(rels) != 1 || rels[0].AddedAt != rel.AddedAt {
		t.Errorf("GetPlaylistMusic failed to match exact relation")
	}

	if err := s.DeletePlaylistMusic(p.PlaylistId, rel.MusicId, MusicSource(rel.Source)); err != nil {
		t.Fatalf("DeletePlaylistMusic failed: %v", err)
	}
}

// ---------------- File Storer Tests ----------------

func TestServerStorage_FileOps(t *testing.T) {
	s, _, cleanup := setupTestServerStorage(t)
	defer cleanup()

	m := Music{MusicId: "f1", Source: YouTubeSource}
	content := []byte("audio_data")

	if err := s.PutMusicFile(m, bytes.NewReader(content)); err != nil {
		t.Fatalf("PutMusicFile failed: %v", err)
	}

	reader, err := s.GetMusicFile(m)
	if err != nil {
		t.Fatalf("GetMusicFile failed: %v", err)
	}

	readData, err := io.ReadAll(reader)

	if closeErr := reader.Close(); closeErr != nil {
		t.Fatalf("Failed closing reader: %v", closeErr)
	}

	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}
	if !bytes.Equal(readData, content) {
		t.Errorf("File content mismatch")
	}

	if err := s.DeleteMusicFile(m); err != nil {
		t.Fatalf("DeleteMusicFile failed: %v", err)
	}
}

// ---------------- Sync Diffing Algorithm Tests ----------------

// ---------------- Sync Diffing Algorithm Tests ----------------

func TestServerStorage_DownloadPlaylistFromServer_DiffLogic(t *testing.T) {
	s, ts, cleanup := setupTestServerStorage(t)
	defer cleanup()

	if err := s.PutUser(UserProfile{UserId: "u1", Token: "tok"}); err != nil {
		t.Fatalf("failed user setup: %v", err)
	}

	localPl, err := s.Storage.PutPlaylist(Playlist{Title: "Diff Test", CoverBlob: []byte{}})
	if err != nil {
		t.Fatalf("failed playlist setup: %v", err)
	}

	_ = s.PutMusic(Music{MusicId: "stale_1", Source: YouTubeSource, Title: "Stale Song", LengthSeconds: 10})
	_ = s.PutMusic(Music{MusicId: "drift_1", Source: YouTubeSource, Title: "Drift Song", LengthSeconds: 10})

	if err := s.Storage.PutPlaylistMusic(PlaylistMusic{
		UserId: "u1", PlaylistId: localPl.PlaylistId, MusicId: "stale_1", Source: int64(YouTubeSource), AddedAt: 100,
	}); err != nil {
		t.Fatalf("failed relation setup: %v", err)
	}

	if err := s.Storage.PutPlaylistMusic(PlaylistMusic{
		UserId: "u1", PlaylistId: localPl.PlaylistId, MusicId: "drift_1", Source: int64(YouTubeSource), AddedAt: 200,
	}); err != nil {
		t.Fatalf("failed relation setup: %v", err)
	}

	ts.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/getPlaylistContent" {
			resp := schemas.GetPlaylistContentResponse{
				Musics: []schemas.Music{
					{MusicId: "drift_1", Source: int64(YouTubeSource), Title: "Kept Song"},
					{MusicId: "new_1", Source: int64(SpotifySource), Title: "New Song"},
				},
				Relations: []schemas.PlaylistMusic{
					{UserId: "u1", PlaylistId: localPl.PlaylistId, MusicId: "drift_1", Source: int64(YouTubeSource), AddedAt: 999},
					{UserId: "u1", PlaylistId: localPl.PlaylistId, MusicId: "new_1", Source: int64(SpotifySource), AddedAt: 500},
				},
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Write([]byte(`{}`))
	})

	// FIX: Added CoverBlob: []byte{} to satisfy NOT NULL constraint during mock remote download
	remotePl := schemas.Playlist{
		UserId:     "u1",
		PlaylistId: localPl.PlaylistId,
		Title:      "Diff Test Remote",
		CoverBlob:  []byte{},
	}
	if err := s.downloadPlaylistFromServer(remotePl, "tok"); err != nil {
		t.Fatalf("downloadPlaylistFromServer failed: %v", err)
	}

	finalRels, err := s.GetAllPlaylistMusic(localPl.PlaylistId)
	if err != nil {
		t.Fatalf("GetPlaylistMusic failed: %v", err)
	}

	if len(finalRels) != 2 {
		t.Fatalf("Diff failed: Expected exactly 2 relations, got %d", len(finalRels))
	}

	hasDrift := false
	hasNew := false
	for _, r := range finalRels {
		if r.MusicId == "stale_1" {
			t.Errorf("Diff failed: stale relation was not deleted")
		}
		if r.MusicId == "drift_1" {
			hasDrift = true
			if r.AddedAt != 999 {
				t.Errorf("Diff failed: drift relation time was not updated to 999, got %d", r.AddedAt)
			}
		}
		if r.MusicId == "new_1" {
			hasNew = true
		}
	}
	if !hasDrift || !hasNew {
		t.Errorf("Diff failed to correctly map and insert new/drifted relations")
	}
}

// ---------------- Offline Degradation Tests ----------------

func TestServerStorage_TryRemote_OfflineDegradation(t *testing.T) {
	s, ts, cleanup := setupTestServerStorage(t)
	defer cleanup()

	if err := s.PutUser(UserProfile{UserId: "u1", Token: "tok"}); err != nil {
		t.Fatalf("failed to put user: %v", err)
	}

	ts.Close()

	// FIX: Added CoverBlob initialization to satisfy NOT NULL constraint
	p := Playlist{Title: "Will trigger offline", CoverBlob: []byte{}}
	if _, err := s.PutPlaylist(p); err != nil {
		t.Fatalf("PutPlaylist should not fail locally when server is down: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if !s.isCurrentlyOffline() {
		t.Error("ServerStorage failed to enter offline mode after a background network failure")
	}
}
