package storages

import (
	"io"
	"log/slog"
	"meowyplayer/handlers"
	"net/http"
	"sync"
)

type ServerStorage struct {
	local      *SQLiteStorage
	endpoints  map[string]string
	offline    bool
	offlineMux sync.Mutex
	syncingMux sync.Mutex
}

func NewServerStorage(local *SQLiteStorage, endpoints map[string]string) *ServerStorage {
	return &ServerStorage{
		local:     local,
		endpoints: endpoints,
		offline:   false,
	}
}

func (s *ServerStorage) getToken() string {
	user, err := s.local.GetUser()
	if err != nil {
		return ""
	}
	return user.Token
}

func (s *ServerStorage) markOffline() {
	s.offlineMux.Lock()
	defer s.offlineMux.Unlock()
	s.offline = true
	slog.Info("enter offline mode")
}

func (s *ServerStorage) isCurrentlyOffline() bool {
	s.offlineMux.Lock()
	defer s.offlineMux.Unlock()
	return s.offline
}

// tryRemote abstracts the preSync -> network call -> offline fallback pattern
func tryRemote[T any, Y any](s *ServerStorage, endpointKey string, req T, resp *Y) {
	if err := s.preSync(); err != nil {
		return
	}
	if err := handlers.SendJSON(s.endpoints[endpointKey], req, resp); err != nil {
		s.markOffline()
	}
}

func (s *ServerStorage) uploadPlaylistToServer(lp Playlist, token string) error {
	// 1. Upload Playlist.
	putReq := handlers.PutPlaylistRequest{
		Token: token,
		Playlist: handlers.Playlist{
			UserId:       lp.UserId,
			PlaylistId:   lp.PlaylistId,
			Title:        lp.Title,
			ModifiedDate: lp.ModifiedDate,
			CoverBlob:    lp.CoverBlob,
		},
	}
	if err := handlers.SendJSON(s.endpoints["putPlaylist"], putReq, &handlers.PutPlaylistResponse{}); err != nil {
		return err
	}

	// 2. Fetch local music.
	musics, err := s.local.GetMusicFromPlaylist(lp.PlaylistId)
	if err != nil || len(musics) == 0 {
		return err
	}

	// 3. Prepare music and relations bulks.
	bulkMusic := make([]handlers.Music, len(musics))
	bulkRelations := make([]handlers.PlaylistMusic, len(musics))

	for i, m := range musics {
		bulkMusic[i] = handlers.Music{
			MusicId:       m.MusicId,
			Source:        handlers.MusicSource(m.Source),
			Title:         m.Title,
			LengthSeconds: m.LengthSeconds,
		}
		bulkRelations[i] = handlers.PlaylistMusic{
			UserId:     lp.UserId,
			PlaylistId: lp.PlaylistId,
			MusicId:    m.MusicId,
			Source:     int64(m.Source),
		}
	}

	// 4. Bulk Uploads
	mReq := handlers.PutMusicBulkRequest{Token: token, Music: bulkMusic}
	if err := handlers.SendJSON(s.endpoints["putMusicBulk"], mReq, &handlers.PutMusicBulkResponse{}); err != nil {
		return err
	}

	linkReq := handlers.PutMusicInPlaylistBulkRequest{Token: token, Relations: bulkRelations}
	return handlers.SendJSON(s.endpoints["putMusicInPlaylistBulk"], linkReq, &handlers.PutMusicInPlaylistBulkResponse{})
}

func (s *ServerStorage) downloadPlaylistFromServer(rp handlers.Playlist, token string) error {
	// 1. Save metadata locally.
	pl := Playlist{
		UserId:       rp.UserId,
		PlaylistId:   rp.PlaylistId,
		Title:        rp.Title,
		ModifiedDate: rp.ModifiedDate,
		CoverBlob:    rp.CoverBlob,
	}
	if _, err := s.local.PutPlaylist(pl); err != nil {
		return err
	}

	// 2. Fetch contents
	reqC := handlers.GetPlaylistContentRequest{Token: token, PlaylistId: rp.PlaylistId}
	var respC handlers.GetPlaylistContentResponse
	if err := handlers.SendJSON(s.endpoints["getPlaylistContent"], reqC, &respC); err != nil {
		return err
	}

	// 3. Save music and relations
	for _, m := range respC.Musics {
		music := Music{
			MusicId:       m.MusicId,
			Source:        MusicSource(m.Source),
			Title:         m.Title,
			LengthSeconds: m.LengthSeconds,
		}
		if err := s.local.PutMusic(music); err != nil {
			return err
		}
		if err := s.local.PutMusicInPlaylist(rp.PlaylistId, m.MusicId, MusicSource(m.Source)); err != nil {
			return err
		}
	}
	return nil
}

func (s *ServerStorage) preSync() error {
	s.syncingMux.Lock()
	defer s.syncingMux.Unlock()

	if !s.isCurrentlyOffline() {
		return nil
	}
	token := s.getToken()

	// 1. Fetch remote playlists
	var resp handlers.GetPlaylistsFromUserResponse
	if err := handlers.SendJSON(s.endpoints["getPlaylistsFromUser"], handlers.GetPlaylistsFromUserRequest{Token: token}, &resp); err != nil {
		return err
	}

	// 2. Fetch and index local playlists
	localPlaylistsSlice, err := s.local.GetPlaylistsFromUser()
	if err != nil {
		return err
	}

	localPlaylists := make(map[int64]Playlist, len(localPlaylistsSlice))
	for _, p := range localPlaylistsSlice {
		localPlaylists[p.PlaylistId] = p
	}

	// 3. Process remote playlists and compare
	for _, rp := range resp.Playlists {
		lp, existsLocally := localPlaylists[rp.PlaylistId]

		if !existsLocally || rp.ModifiedDate > lp.ModifiedDate {
			// Server is newer or local is missing -> Download
			if err := s.downloadPlaylistFromServer(rp, token); err != nil {
				return err
			}
		} else if lp.ModifiedDate > rp.ModifiedDate {
			// Local is newer -> Upload
			if err := s.uploadPlaylistToServer(lp, token); err != nil {
				return err
			}
		}

		// Remove from map: it exists remotely, so it's not a local-only playlist
		delete(localPlaylists, rp.PlaylistId)
	}

	// 4. Anything remaining in localPlaylists is local-only. Upload them.
	for _, lp := range localPlaylists {
		if err := s.uploadPlaylistToServer(lp, token); err != nil {
			return err
		}
	}

	s.offlineMux.Lock()
	s.offline = false
	slog.Info("enter online mode")
	s.offlineMux.Unlock()

	return nil
}

func (s *ServerStorage) PutUser(userProfile UserProfile) error {
	return s.local.PutUser(userProfile)
}

func (s *ServerStorage) GetUser() (UserProfile, error) {
	return s.local.GetUser()
}

func (s *ServerStorage) DeleteUser() error {
	return s.local.DeleteUser()
}

func (s *ServerStorage) PutPlaylist(playlist Playlist) (Playlist, error) {
	res, err := s.local.PutPlaylist(playlist)
	if err != nil {
		return Playlist{}, err
	}

	req := handlers.PutPlaylistRequest{
		Token: s.getToken(),
		Playlist: handlers.Playlist{
			UserId: res.UserId, PlaylistId: res.PlaylistId,
			Title: res.Title, ModifiedDate: res.ModifiedDate, CoverBlob: res.CoverBlob,
		},
	}
	tryRemote(s, "putPlaylist", req, &handlers.PutPlaylistResponse{})
	return res, nil
}

func (s *ServerStorage) GetPlaylist(playlistId int64) (Playlist, error) {
	return s.local.GetPlaylist(playlistId)
}

func (s *ServerStorage) DeletePlaylist(playlistId int64) error {
	if err := s.local.DeletePlaylist(playlistId); err != nil {
		return err
	}
	tryRemote(s, "deletePlaylist", handlers.DeletePlaylistRequest{Token: s.getToken(), PlaylistId: playlistId}, &handlers.DeletePlaylistResponse{})
	return nil
}

func (s *ServerStorage) PutMusic(music Music) error {
	if err := s.local.PutMusic(music); err != nil {
		return err
	}
	req := handlers.PutMusicRequest{
		Token: s.getToken(),
		Music: handlers.Music{
			MusicId: music.MusicId, Source: handlers.MusicSource(music.Source),
			Title: music.Title, LengthSeconds: music.LengthSeconds,
		},
	}
	tryRemote(s, "putMusic", req, &handlers.PutMusicResponse{})
	return nil
}

func (s *ServerStorage) GetMusic(musicId string, source MusicSource) (Music, error) {
	return s.local.GetMusic(musicId, source)
}

func (s *ServerStorage) DeleteMusic(musicId string, source MusicSource) error {
	return s.local.DeleteMusic(musicId, source)
}

func (s *ServerStorage) PutMusicFile(music Music, content io.Reader) error {
	return s.local.PutMusicFile(music, content)
}

func (s *ServerStorage) GetMusicFile(music Music) (io.ReadSeekCloser, error) {
	return s.local.GetMusicFile(music)
}

func (s *ServerStorage) DeleteMusicFile(music Music) error {
	return s.local.DeleteMusicFile(music)
}

func (s *ServerStorage) GetPlaylistsFromUser() ([]Playlist, error) {
	return s.local.GetPlaylistsFromUser()
}

func (s *ServerStorage) GetMusicFromPlaylist(playlistId int64) ([]Music, error) {
	return s.local.GetMusicFromPlaylist(playlistId)
}

func (s *ServerStorage) PutMusicInPlaylist(playlistId int64, musicId string, source MusicSource) error {
	if err := s.local.PutMusicInPlaylist(playlistId, musicId, source); err != nil {
		return err
	}
	req := handlers.PutMusicInPlaylistRequest{
		Token: s.getToken(), PlaylistId: playlistId, MusicId: musicId, Source: handlers.MusicSource(source),
	}
	tryRemote(s, "putMusicInPlaylist", req, &handlers.PutMusicInPlaylistResponse{})
	return nil
}

func (s *ServerStorage) DeleteMusicFromPlaylist(playlistId int64, musicId string, source MusicSource) error {
	if err := s.local.DeleteMusicFromPlaylist(playlistId, musicId, source); err != nil {
		return err
	}
	req := handlers.DeleteMusicFromPlaylistRequest{
		Token: s.getToken(), PlaylistId: playlistId, MusicId: musicId, Source: handlers.MusicSource(source),
	}
	tryRemote(s, "deleteMusicFromPlaylist", req, &handlers.DeleteMusicFromPlaylistResponse{})
	return nil
}

func (s *ServerStorage) Close() error {
	http.DefaultClient.CloseIdleConnections()
	return s.local.Close()
}
