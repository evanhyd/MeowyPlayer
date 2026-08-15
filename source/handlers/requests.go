package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"meowyplayer/storages"
	"net/http"
)

// Request
type RegisterRequest struct {
	UserId   string            `json:"user_id"`
	Username string            `json:"username"`
	Language storages.Language `json:"language"`
	Password string            `json:"password"`
}

type LoginRequest struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	Token string `json:"token"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type MeRequest struct {
	Token string `json:"token"`
}

// Response
type ErrorResponse struct {
	Error string `json:"error"`
}

type RegisterResponse struct {
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RefreshResponse struct {
	Token string `json:"token"`
}

type ResetPasswordResponse struct {
}

type MeResponse struct {
	Profile storages.UserProfile `json:"profile"`
}

func SendJSON[T any, Y any](url string, request T, response *Y) error {
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

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	// Check if it is error response.
	if rsp.StatusCode >= 200 && rsp.StatusCode <= 299 {
		json.NewDecoder(rsp.Body).Decode(response)
		return nil
	} else {
		errorRsp := ErrorResponse{}
		json.NewDecoder(rsp.Body).Decode(&errorRsp)
		return errors.New(errorRsp.Error)
	}
}
