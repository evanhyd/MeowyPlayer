package storages

import (
	"encoding/json"
	"meowyplayer/handlers"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// setupTestEnvironment initializes a mock HTTP server, an in-memory SQLite DB, and the ServerStorage.
func setupTestEnvironment(t *testing.T, handler http.HandlerFunc) (*ServerStorage, func()) {
	ts := httptest.NewServer(handler)

	// Create temporary directory for music files
	tempDir, err := os.MkdirTemp("", "meowyplayer_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Use SQLite's special :memory: identifier for an in-memory database
	dbPath := ":memory:"
	musicPath := filepath.Join(tempDir, "music")

	localStore := NewSQLiteStorage(dbPath, musicPath)
	if localStore == nil {
		t.Fatalf("Failed to initialize SQLiteStorage")
	}

	// Seed a test user so getToken() has something to return
	err = localStore.PutUser(UserProfile{
		UserId:   "test_user_123",
		Username: "TestUser",
		Token:    "mock_token_abc",
	})
	if err != nil {
		t.Fatalf("Failed to seed test user: %v", err)
	}

	endpoints := map[string]string{
		"putPlaylist":             ts.URL + "/putPlaylist",
		"getPlaylistsFromUser":    ts.URL + "/getPlaylistsFromUser",
		"getPlaylistContent":      ts.URL + "/getPlaylistContent",
		"putMusicInPlaylist":      ts.URL + "/putMusicInPlaylist",
		"deletePlaylist":          ts.URL + "/deletePlaylist",
		"putMusic":                ts.URL + "/putMusic",
		"deleteMusicFromPlaylist": ts.URL + "/deleteMusicFromPlaylist",
		"putMusicBulk":            ts.URL + "/putMusicBulk",
		"putMusicInPlaylistBulk":  ts.URL + "/putMusicInPlaylistBulk",
	}

	serverStore := NewServerStorage(localStore, endpoints)

	cleanup := func() {
		ts.Close()
		serverStore.Close()
		os.RemoveAll(tempDir)
	}

	return serverStore, cleanup
}

// ---------------- Playlist Tests ----------------

func TestServerStorage_PutPlaylist(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(handlers.PutPlaylistResponse{})
	}
	store, cleanup := setupTestEnvironment(t, handler)
	defer cleanup()

	p := Playlist{
		Title:     "My Awesome Playlist",
		CoverBlob: []byte("mock_cover_image_data"),
	}

	savedPlaylist, err := store.PutPlaylist(p)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 1. Verify local storage succeeded
	localPlaylist, err := store.local.GetPlaylist(savedPlaylist.PlaylistId)
	if err != nil {
		t.Fatalf("Expected playlist to exist locally, got error: %v", err)
	}

	// 2. Verify CoverBlob is not nil and matches
	if localPlaylist.CoverBlob == nil || string(localPlaylist.CoverBlob) != "mock_cover_image_data" {
		t.Errorf("Expected CoverBlob to be 'mock_cover_image_data', got %v", localPlaylist.CoverBlob)
	}

	// 3. Verify we did not go offline (synchronous check)
	if store.isCurrentlyOffline() {
		t.Errorf("Expected store to remain online")
	}
}

func TestServerStorage_DeletePlaylist(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(handlers.DeletePlaylistResponse{})
	}
	store, cleanup := setupTestEnvironment(t, handler)
	defer cleanup()

	pl, _ := store.PutPlaylist(Playlist{Title: "To Delete", CoverBlob: []byte("cover")})

	err := store.DeletePlaylist(pl.PlaylistId)
	if err != nil {
		t.Fatalf("Expected DeletePlaylist to succeed, got %v", err)
	}

	// Verify local deletion
	_, err = store.local.GetPlaylist(pl.PlaylistId)
	if err == nil {
		t.Errorf("Expected playlist to be deleted locally")
	}
}

// ---------------- Music Tests ----------------

func TestServerStorage_PutMusic(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(handlers.PutMusicResponse{})
	}
	store, cleanup := setupTestEnvironment(t, handler)
	defer cleanup()

	m := Music{
		MusicId:       "song_123",
		Source:        YouTubeSource,
		Title:         "Never Gonna Give You Up",
		LengthSeconds: 212,
	}

	err := store.PutMusic(m)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	localM, err := store.local.GetMusic("song_123", YouTubeSource)
	if err != nil {
		t.Fatalf("Expected music to exist locally, got %v", err)
	}
	if localM.Title != "Never Gonna Give You Up" {
		t.Errorf("Title mismatch: got %v", localM.Title)
	}
}

// ---------------- Playlist <-> Music Relation Tests ----------------

func TestServerStorage_PutMusicInPlaylist(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(handlers.PutMusicInPlaylistResponse{})
	}
	store, cleanup := setupTestEnvironment(t, handler)
	defer cleanup()

	pl, _ := store.local.PutPlaylist(Playlist{Title: "My Mix", CoverBlob: []byte("mix_cover")})
	store.local.PutMusic(Music{MusicId: "song_abc", Source: SpotifySource, Title: "Track 1"})

	err := store.PutMusicInPlaylist(pl.PlaylistId, "song_abc", SpotifySource)
	if err != nil {
		t.Fatalf("Expected PutMusicInPlaylist to succeed, got %v", err)
	}

	musics, err := store.local.GetMusicFromPlaylist(pl.PlaylistId)
	if err != nil || len(musics) != 1 {
		t.Fatalf("Expected 1 song in playlist, got %v (err: %v)", len(musics), err)
	}
}

func TestServerStorage_DeleteMusicFromPlaylist(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(handlers.DeleteMusicFromPlaylistResponse{})
	}
	store, cleanup := setupTestEnvironment(t, handler)
	defer cleanup()

	pl, _ := store.local.PutPlaylist(Playlist{Title: "My Mix", CoverBlob: []byte("mix_cover")})
	store.local.PutMusic(Music{MusicId: "song_xyz", Source: SpotifySource})
	store.local.PutMusicInPlaylist(pl.PlaylistId, "song_xyz", SpotifySource)

	err := store.DeleteMusicFromPlaylist(pl.PlaylistId, "song_xyz", SpotifySource)
	if err != nil {
		t.Fatalf("Expected DeleteMusicFromPlaylist to succeed, got %v", err)
	}

	musics, _ := store.local.GetMusicFromPlaylist(pl.PlaylistId)
	if len(musics) != 0 {
		t.Errorf("Expected playlist to be empty, found %d songs", len(musics))
	}
}

// ---------------- Offline / PreSync Tests ----------------

func TestServerStorage_NetworkFailure_SetsOffline(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Error: "server crash"})
	}
	store, cleanup := setupTestEnvironment(t, handler)
	defer cleanup()

	p := Playlist{Title: "Offline Playlist", CoverBlob: []byte("test_blob")}

	// Triggers synchronous network call which will hit the 500 error
	store.PutPlaylist(p)

	if !store.isCurrentlyOffline() {
		t.Errorf("Expected store to be marked offline after 500 error")
	}
}

func TestServerStorage_PreSync_Recovery_FullState(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		// Mock responses for downloading a cloud playlist
		case "/getPlaylistsFromUser":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(handlers.GetPlaylistsFromUserResponse{
				Playlists: []handlers.Playlist{
					{PlaylistId: 999, Title: "Cloud Sync Playlist", CoverBlob: []byte("synced_cover_blob"), ModifiedDate: 999999999},
				},
			})
		case "/getPlaylistContent":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(handlers.GetPlaylistContentResponse{
				Musics: []handlers.Music{
					{MusicId: "song_1", Title: "Cloud Song", Source: int64(YouTubeSource)},
				},
			})

		// Mock responses for uploading a local playlist
		case "/putPlaylist":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(handlers.PutPlaylistResponse{})
		case "/putMusicBulk":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(handlers.PutMusicBulkResponse{})
		case "/putMusicInPlaylistBulk":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(handlers.PutMusicInPlaylistBulkResponse{})
		case "/putMusicInPlaylist":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(handlers.PutMusicInPlaylistResponse{})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}

	store, cleanup := setupTestEnvironment(t, handler)
	defer cleanup()

	// Seed local data (this will get uploaded during preSync because it doesn't exist remotely)
	localPl, _ := store.local.PutPlaylist(Playlist{Title: "Local Target", CoverBlob: []byte("local_blob")})
	store.local.PutMusic(Music{MusicId: "song_new", Source: YouTubeSource, Title: "New Song"})

	store.markOffline()

	// Trigger modifying action -> invokes synchronous preSync()
	err := store.PutMusicInPlaylist(localPl.PlaylistId, "song_new", YouTubeSource)
	if err != nil {
		t.Fatalf("Expected PutMusicInPlaylist to succeed locally, got %v", err)
	}

	// 1. Verify cloud playlist 999 was pulled down
	syncedPlaylist, err := store.local.GetPlaylist(999)
	if err != nil {
		t.Fatalf("Expected to find synced playlist 999 locally, got %v", err)
	}
	if syncedPlaylist.CoverBlob == nil || string(syncedPlaylist.CoverBlob) != "synced_cover_blob" {
		t.Errorf("Expected synced CoverBlob 'synced_cover_blob', got %v", syncedPlaylist.CoverBlob)
	}

	// 2. Verify cloud music was synced locally
	syncedMusic, _ := store.local.GetMusicFromPlaylist(999)
	if len(syncedMusic) == 0 || syncedMusic[0].MusicId != "song_1" {
		t.Errorf("Expected song_1 to be synced to playlist 999")
	}

	// 3. Verify the store is marked back ONLINE
	if store.isCurrentlyOffline() {
		t.Errorf("Expected store to be back online after a successful preSync")
	}
}

func TestServerStorage_PreSync_Fails_ContinuesOffline(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError) // Server is still down
	}
	store, cleanup := setupTestEnvironment(t, handler)
	defer cleanup()

	store.local.PutMusic(Music{MusicId: "m_fail_test", Source: YouTubeSource})
	store.markOffline()

	// Operation happens locally, preSync fails silently but doesn't block the local save
	err := store.PutMusic(Music{MusicId: "m_fail_test", Title: "Updated Offline", Source: YouTubeSource})
	if err != nil {
		t.Fatalf("Expected local operation to succeed even if preSync fails, got %v", err)
	}

	// Verify local operation succeeded
	updatedM, _ := store.local.GetMusic("m_fail_test", YouTubeSource)
	if updatedM.Title != "Updated Offline" {
		t.Errorf("Expected local update to persist")
	}

	// Verify we are still offline
	if !store.isCurrentlyOffline() {
		t.Errorf("Expected store to remain offline because preSync failed")
	}
}
