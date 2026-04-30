package mutil

import "strconv"

func PlaylistIdToString(playlistId int64) string {
	return strconv.FormatInt(playlistId, 16)
}
