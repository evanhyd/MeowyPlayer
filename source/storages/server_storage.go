package storages

import (
	"fmt"
	"log/slog"
	"meowyplayer/schemas"
	"net/http"
	"sync"
	"time"
)

var _ Storage = (*ServerStorage)(nil)

type ServerStorage struct {
	Storage

	httpClient       *http.Client
	endpointProvider EndpointProvider
	offline          bool
	syncedOnce       bool
	offlineMux       sync.Mutex
	syncingMux       sync.Mutex
	stopWatch        chan struct{}
}

func NewServerStorage(client *http.Client, localStorage Storage, provider EndpointProvider) *ServerStorage {
	return &ServerStorage{
		Storage:          localStorage,
		httpClient:       client,
		endpointProvider: provider,
		offline:          false,
		syncedOnce:       false,
		stopWatch:        make(chan struct{}),
	}
}

func (s *ServerStorage) getToken() string {
	user, err := s.Storage.GetUser() // Call embedded method explicitly
	if err != nil {
		return ""
	}
	return user.Token
}

func (s *ServerStorage) markOffline() {
	s.offlineMux.Lock()
	if s.offline {
		s.offlineMux.Unlock()
		return
	}
	s.offline = true
	slog.Info("entered offline mode; starting background reconnection watcher")
	s.offlineMux.Unlock()

	go s.watchForReconnection()
}

func (s *ServerStorage) markOnline() {
	s.offlineMux.Lock()
	defer s.offlineMux.Unlock()
	s.offline = false
	slog.Info("entered online mode")
}

func (s *ServerStorage) isCurrentlyOffline() bool {
	s.offlineMux.Lock()
	defer s.offlineMux.Unlock()
	return s.offline
}

func (s *ServerStorage) hasSynced() bool {
	s.offlineMux.Lock()
	defer s.offlineMux.Unlock()
	return s.syncedOnce
}

func (s *ServerStorage) setSynced(val bool) {
	s.offlineMux.Lock()
	defer s.offlineMux.Unlock()
	s.syncedOnce = val
}

func (s *ServerStorage) checkOnline() bool {
	token := s.getToken()
	if token == "" {
		return false
	}

	endpoint, err := s.endpointProvider.GetEndpoint("me")
	if err != nil {
		return false
	}

	request := schemas.MeRequest{Token: token}
	response := schemas.MeResponse{}
	if err := schemas.SendJSON(s.httpClient, endpoint, request, &response); err != nil {
		return false
	}

	if err := s.Storage.PutUser(UserProfile{
		UserId:           response.UserId,
		Username:         response.Username,
		Language:         response.Language,
		RegistrationDate: response.RegistrationDate,
		Token:            token,
	}); err != nil {
		slog.Error("failed to update user profile during online check", "error", err)
	}

	return true
}

func (s *ServerStorage) watchForReconnection() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopWatch:
			return
		case <-ticker.C:
			if !s.isCurrentlyOffline() {
				return
			}

			if s.checkOnline() {
				slog.Info("connection restored to server")
				s.markOnline()

				go func() {
					if err := s.preSync(); err != nil {
						slog.Error("error during catch-up preSync", "error", err)
					}
				}()
				return
			}
		}
	}
}

func tryRemote[T any, Y any](s *ServerStorage, endpointKey string, req T, resp *Y) {
	go func() {
		if s.isCurrentlyOffline() {
			return
		}

		endpoint, err := s.endpointProvider.GetEndpoint(endpointKey)
		if err != nil {
			return
		}

		if err := schemas.SendJSON(s.httpClient, endpoint, req, resp); err != nil {
			slog.Warn("remote request failed, switching to offline mode", "endpoint", endpointKey, "error", err)
			s.markOffline()
		}
	}()
}

func (s *ServerStorage) uploadPlaylistToServer(lp Playlist, token string) error {
	epPutPlaylist, err := s.endpointProvider.GetEndpoint("putPlaylist")
	if err != nil {
		return err
	}
	epPutMusicBulk, err := s.endpointProvider.GetEndpoint("putMusicBulk")
	if err != nil {
		return err
	}
	epPutMusicInPlaylistBulk, err := s.endpointProvider.GetEndpoint("putMusicInPlaylistBulk")
	if err != nil {
		return err
	}

	putReq := schemas.PutPlaylistRequest{
		Token: token,
		Playlist: schemas.Playlist{
			UserId:       lp.UserId,
			PlaylistId:   lp.PlaylistId,
			Title:        lp.Title,
			ModifiedDate: lp.ModifiedDate,
			CoverBlob:    lp.CoverBlob,
		},
	}

	localRelations, err := s.Storage.GetAllPlaylistMusic(lp.PlaylistId)
	if err != nil {
		return err
	}

	if len(localRelations) == 0 {
		return schemas.SendJSON(s.httpClient, epPutPlaylist, putReq, &schemas.PutPlaylistResponse{})
	}

	bulkMusic := make([]schemas.Music, len(localRelations))
	bulkRelations := make([]schemas.PlaylistMusic, len(localRelations))

	for i, rel := range localRelations {
		m, err := s.Storage.GetMusic(rel.MusicId, MusicSource(rel.Source))
		if err != nil {
			return fmt.Errorf("failed fetching music metadata for sync: %v", err)
		}

		bulkMusic[i] = schemas.Music{
			MusicId:       m.MusicId,
			Source:        schemas.MusicSource(m.Source),
			Title:         m.Title,
			LengthSeconds: m.LengthSeconds,
		}
		bulkRelations[i] = schemas.PlaylistMusic{
			UserId:     lp.UserId,
			PlaylistId: lp.PlaylistId,
			MusicId:    rel.MusicId,
			Source:     rel.Source,
			AddedAt:    rel.AddedAt,
		}
	}

	mReq := schemas.PutMusicBulkRequest{Token: token, Music: bulkMusic}
	linkReq := schemas.PutMusicInPlaylistBulkRequest{Token: token, Relations: bulkRelations}

	errChan := make(chan error, 2)

	go func() {
		errChan <- schemas.SendJSON(s.httpClient, epPutPlaylist, putReq, &schemas.PutPlaylistResponse{})
	}()

	go func() {
		errChan <- schemas.SendJSON(s.httpClient, epPutMusicBulk, mReq, &schemas.PutMusicBulkResponse{})
	}()

	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return schemas.SendJSON(s.httpClient, epPutMusicInPlaylistBulk, linkReq, &schemas.PutMusicInPlaylistBulkResponse{})
}

func (s *ServerStorage) downloadPlaylistFromServer(rp schemas.Playlist, token string) error {
	epGetPlaylistContent, err := s.endpointProvider.GetEndpoint("getPlaylistContent")
	if err != nil {
		return err
	}

	pl := Playlist{
		UserId:       rp.UserId,
		PlaylistId:   rp.PlaylistId,
		Title:        rp.Title,
		ModifiedDate: rp.ModifiedDate,
		CoverBlob:    rp.CoverBlob,
	}
	if pl.CoverBlob == nil {
		pl.CoverBlob = make([]byte, 0)
	}

	if _, err := s.Storage.PutPlaylist(pl); err != nil {
		return err
	}

	reqC := schemas.GetPlaylistContentRequest{Token: token, PlaylistId: rp.PlaylistId}
	var respC schemas.GetPlaylistContentResponse
	if err := schemas.SendJSON(s.httpClient, epGetPlaylistContent, reqC, &respC); err != nil {
		return err
	}

	for _, m := range respC.Musics {
		music := Music{
			MusicId:       m.MusicId,
			Source:        MusicSource(m.Source),
			Title:         m.Title,
			LengthSeconds: m.LengthSeconds,
		}
		if err := s.Storage.PutMusic(music); err != nil {
			return err
		}
	}

	localRels, err := s.Storage.GetAllPlaylistMusic(rp.PlaylistId)
	if err != nil {
		return err
	}

	type musicKey struct {
		Source  int64
		MusicId string
	}

	localMap := make(map[musicKey]PlaylistMusic, len(localRels))
	for _, lr := range localRels {
		localMap[musicKey{Source: lr.Source, MusicId: lr.MusicId}] = lr
	}

	serverMap := make(map[musicKey]schemas.PlaylistMusic, len(respC.Relations))
	for _, rr := range respC.Relations {
		serverMap[musicKey{Source: rr.Source, MusicId: rr.MusicId}] = rr
	}

	for key, lr := range localMap {
		if _, exists := serverMap[key]; !exists {
			if err := s.Storage.DeletePlaylistMusic(rp.PlaylistId, lr.MusicId, MusicSource(lr.Source)); err != nil {
				return err
			}
		}
	}

	for key, rr := range serverMap {
		lr, exists := localMap[key]
		if !exists || lr.AddedAt != rr.AddedAt {
			if err := s.Storage.PutPlaylistMusic(PlaylistMusic{
				UserId:     rr.UserId,
				PlaylistId: rr.PlaylistId,
				MusicId:    rr.MusicId,
				Source:     rr.Source,
				AddedAt:    rr.AddedAt,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *ServerStorage) preSync() error {
	if s.isCurrentlyOffline() {
		return nil
	}

	s.syncingMux.Lock()
	defer s.syncingMux.Unlock()

	token := s.getToken()
	if token == "" {
		return nil
	}

	endpoint, err := s.endpointProvider.GetEndpoint("getPlaylistsFromUser")
	if err != nil {
		return nil
	}

	var resp schemas.GetPlaylistsFromUserResponse
	if err := schemas.SendJSON(s.httpClient, endpoint, schemas.GetPlaylistsFromUserRequest{Token: token}, &resp); err != nil {
		slog.Warn("preSync failed to reach server, going offline", "error", err)
		s.markOffline()
		return nil
	}

	localPlaylistsSlice, err := s.Storage.GetPlaylists()
	if err != nil {
		return err
	}

	localPlaylists := make(map[int64]Playlist, len(localPlaylistsSlice))
	for _, p := range localPlaylistsSlice {
		localPlaylists[p.PlaylistId] = p
	}

	for _, rp := range resp.Playlists {
		lp, existsLocally := localPlaylists[rp.PlaylistId]
		if !existsLocally || rp.ModifiedDate > lp.ModifiedDate {
			if err := s.downloadPlaylistFromServer(rp, token); err != nil {
				slog.Error("failed downloading playlist from server", "playlistId", rp.PlaylistId, "error", err)
			}
		} else if lp.ModifiedDate > rp.ModifiedDate {
			if err := s.uploadPlaylistToServer(lp, token); err != nil {
				slog.Error("failed uploading playlist to server", "playlistId", lp.PlaylistId, "error", err)
			}
		}

		delete(localPlaylists, rp.PlaylistId)
	}

	for _, lp := range localPlaylists {
		if err := s.uploadPlaylistToServer(lp, token); err != nil {
			slog.Error("failed uploading local-only playlist to server", "playlistId", lp.PlaylistId, "error", err)
		}
	}

	s.setSynced(true)
	s.markOnline()
	return nil
}

// ---------------- Overridden User Storer Methods ----------------

func (s *ServerStorage) DeleteUser() error {
	s.setSynced(false)
	return s.Storage.DeleteUser()
}

// ---------------- Overridden Playlist Storer Methods ----------------

func (s *ServerStorage) PutPlaylist(playlist Playlist) (Playlist, error) {
	res, err := s.Storage.PutPlaylist(playlist)
	if err != nil {
		return Playlist{}, err
	}

	req := schemas.PutPlaylistRequest{
		Token: s.getToken(),
		Playlist: schemas.Playlist{
			UserId:       res.UserId,
			PlaylistId:   res.PlaylistId,
			Title:        res.Title,
			ModifiedDate: res.ModifiedDate,
			CoverBlob:    res.CoverBlob,
		},
	}
	tryRemote(s, "putPlaylist", req, &schemas.PutPlaylistResponse{})
	return res, nil
}

func (s *ServerStorage) GetPlaylists() ([]Playlist, error) {
	if !s.hasSynced() && !s.isCurrentlyOffline() {
		if err := s.preSync(); err != nil {
			slog.Error("preSync error in GetPlaylists", "error", err)
		}
	}
	return s.Storage.GetPlaylists()
}

func (s *ServerStorage) DeletePlaylist(playlistId int64) error {
	if err := s.Storage.DeletePlaylist(playlistId); err != nil {
		return err
	}
	tryRemote(s, "deletePlaylist", schemas.DeletePlaylistRequest{Token: s.getToken(), PlaylistId: playlistId}, &schemas.DeletePlaylistResponse{})
	return nil
}

func (s *ServerStorage) PutPlaylistMusic(relation PlaylistMusic) error {
	if err := s.Storage.PutPlaylistMusic(relation); err != nil {
		return err
	}
	req := schemas.PutMusicInPlaylistRequest{
		Token:      s.getToken(),
		PlaylistId: relation.PlaylistId,
		MusicId:    relation.MusicId,
		Source:     schemas.MusicSource(relation.Source),
		AddedAt:    relation.AddedAt,
	}
	tryRemote(s, "putMusicInPlaylist", req, &schemas.PutMusicInPlaylistResponse{})
	return nil
}

func (s *ServerStorage) DeletePlaylistMusic(playlistId int64, musicId string, source MusicSource) error {
	if err := s.Storage.DeletePlaylistMusic(playlistId, musicId, source); err != nil {
		return err
	}
	req := schemas.DeleteMusicFromPlaylistRequest{
		Token:      s.getToken(),
		PlaylistId: playlistId,
		MusicId:    musicId,
		Source:     schemas.MusicSource(source),
	}
	tryRemote(s, "deleteMusicFromPlaylist", req, &schemas.DeleteMusicFromPlaylistResponse{})
	return nil
}

// ---------------- Overridden Music Storer Methods ----------------

func (s *ServerStorage) PutMusic(music Music) error {
	if err := s.Storage.PutMusic(music); err != nil {
		return err
	}
	req := schemas.PutMusicRequest{
		Token: s.getToken(),
		Music: schemas.Music{
			MusicId:       music.MusicId,
			Source:        schemas.MusicSource(music.Source),
			Title:         music.Title,
			LengthSeconds: music.LengthSeconds,
		},
	}
	tryRemote(s, "putMusic", req, &schemas.PutMusicResponse{})
	return nil
}

// ---------------- Overridden File Storer Methods ----------------

func (s *ServerStorage) Close() error {
	close(s.stopWatch)
	s.httpClient.CloseIdleConnections()
	return s.Storage.Close()
}
