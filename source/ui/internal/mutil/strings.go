package mutil

import (
	"fmt"
	"strconv"
)

func PlaylistIdToString(playlistId int64) string {
	return strconv.FormatInt(playlistId, 16)
}

func SecondsToTime(seconds int64) string {
	mins := seconds / 60
	secs := seconds - mins*60
	return fmt.Sprintf("%02d:%02d", mins, secs)
}
