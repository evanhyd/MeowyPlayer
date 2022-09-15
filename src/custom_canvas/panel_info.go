package custom_canvas

type PanelInfo struct {
	AlbumSearchList   *SearchList[AlbumInfo]
	MusicSearchList   *SearchList[MusicInfo]
	SelectedAlbumInfo *AlbumInfo
}

func NewPanelInfo() *PanelInfo {
	return &PanelInfo{AlbumSearchList: nil, MusicSearchList: nil, SelectedAlbumInfo: nil}
}
