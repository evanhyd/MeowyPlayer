package schemas

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type MusicSource = int64
type Language = int64

const (
	UnknownSource MusicSource = iota
	YouTubeSource
	SpotifySource
)

const (
	LangEnglish Language = iota
	LangFrench
	LangChinese
	LangJapanese
)

type UserProfile struct {
	UserId           string   `json:"userId"`
	Username         string   `json:"username"`
	Language         Language `json:"language"`
	RegistrationDate int64    `json:"registrationDate"`
	Token            string   `json:"token"`
}

type Playlist struct {
	UserId       string `json:"userId"`
	PlaylistId   int64  `json:"playlistId"`
	Title        string `json:"title"`
	ModifiedDate int64  `json:"modifiedDate"`
	CoverBlob    []byte `json:"coverBlob"`
}

type Music struct {
	MusicId       string      `json:"musicId"`
	Source        MusicSource `json:"source"`
	Title         string      `json:"title"`
	LengthSeconds int64       `json:"lengthSeconds"`
}

type PlaylistMusic struct {
	UserId       string `json:"userId"`
	PlaylistId   int64  `json:"playlistId"`
	MusicId      string `json:"musicId"`
	Source       int64  `json:"source"`
	ModifiedDate int64  `json:"modifiedDate"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Register.
type RegisterRequest struct {
	UserId   string   `json:"userId"`
	Username string   `json:"username"`
	Language Language `json:"language"`
	Password string   `json:"password"`
}

type RegisterResponse struct {
}

// Login.
type LoginRequest struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Refresh.
type RefreshRequest struct {
	Token string `json:"token"`
}

type RefreshResponse struct {
	Token string `json:"token"`
}

// ResetPassword
type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type ResetPasswordResponse struct {
}

// Me.
type MeRequest struct {
	Token string `json:"token"`
}

type MeResponse struct {
	UserId           string   `json:"userId"`
	Username         string   `json:"username"`
	Language         Language `json:"language"`
	RegistrationDate int64    `json:"registrationDate"` // Unix
}

// GetPlaylist (Lightweight: Returns only the playlist metadata)
type GetPlaylistRequest struct {
	Token      string `json:"token"`
	PlaylistId int64  `json:"playlistId"`
}
type GetPlaylistResponse struct {
	Playlist Playlist `json:"playlist"`
}

// GetPlaylistContent (Heavyweight: Returns playlist, musics, and relations)
type GetPlaylistContentRequest struct {
	Token      string `json:"token"`
	PlaylistId int64  `json:"playlistId"`
}
type GetPlaylistContentResponse struct {
	Playlist  Playlist        `json:"playlist"`
	Musics    []Music         `json:"musics"`
	Relations []PlaylistMusic `json:"relations"`
}

// PutPlaylist
type PutPlaylistRequest struct {
	Token    string   `json:"token"`
	Playlist Playlist `json:"playlist"`
}
type PutPlaylistResponse struct {
	Playlist Playlist `json:"playlist"`
}

// DeletePlaylist
type DeletePlaylistRequest struct {
	Token      string `json:"token"`
	PlaylistId int64  `json:"playlistId"`
}
type DeletePlaylistResponse struct {
}

// GetMusic
type GetMusicRequest struct {
	Token   string      `json:"token"`
	MusicId string      `json:"musicId"`
	Source  MusicSource `json:"source"`
}
type GetMusicResponse struct {
	Music Music `json:"music"`
}

// PutMusic
type PutMusicRequest struct {
	Token string `json:"token"`
	Music Music  `json:"music"`
}
type PutMusicResponse struct {
}

// PutMusicBulk
type PutMusicBulkRequest struct {
	Token string  `json:"token"`
	Music []Music `json:"music"`
}
type PutMusicBulkResponse struct {
}

// GetPlaylists
type GetPlaylistsRequest struct {
	Token string `json:"token"`
}
type GetPlaylistsResponse struct {
	Playlists []Playlist `json:"playlists"`
}

// GetAllMusic
type GetAllMusicRequest struct {
	Token      string `json:"token"`
	PlaylistId int64  `json:"playlistId"`
}
type GetAllMusicResponse struct {
	Musics []Music `json:"musics"`
}

// PutPlaylistMusic
type PutPlaylistMusicRequest struct {
	Token         string        `json:"token"`
	PlaylistMusic PlaylistMusic `json:"playlistMusic"`
}
type PutPlaylistMusicResponse struct {
}

// PutPlaylistMusicBulk
type PutPlaylistMusicBulkRequest struct {
	Token         string          `json:"token"`
	PlaylistMusic []PlaylistMusic `json:"playlistMusic"`
}
type PutPlaylistMusicBulkResponse struct {
}

// DeleteMusicFromPlaylist
type DeletePlaylistMusicRequest struct {
	Token         string        `json:"token"`
	PlaylistMusic PlaylistMusic `json:"playlistMusic"`
}
type DeletePlaylistMusicResponse struct {
}

func SendJSON[T any, Y any](client *http.Client, url string, request T, response *Y) error {
	// Encode request object.
	buffer := bytes.Buffer{}
	err := json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		return err
	}

	// Send request.
	req, err := http.NewRequest("POST", url, &buffer)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	// Check if it is error response.
	if rsp.StatusCode >= 200 && rsp.StatusCode <= 299 {
		return json.NewDecoder(rsp.Body).Decode(response)
	} else {
		errorRsp := ErrorResponse{}
		if err := json.NewDecoder(rsp.Body).Decode(&errorRsp); err != nil || errorRsp.Error == "" {
			return fmt.Errorf("server error http %d", rsp.StatusCode)
		}
		return errors.New(errorRsp.Error)
	}
}

func ReplyJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func ReplyError(w http.ResponseWriter, code int, msg string) {
	ReplyJSON(w, code, ErrorResponse{Error: msg})
}
