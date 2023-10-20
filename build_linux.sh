# go install fyne.io/fyne/v2/cmd/fyne@latest
FYNE_TOOL_PATH=$(go env GOPATH)

# Icon
${FYNE_TOOL_PATH}/bin/fyne bundle -o source/resource/Icon.go --package resource --name MissingIcon asset/missing_asset.png
${FYNE_TOOL_PATH}/bin/fyne bundle -o source/resource/Icon.go --append --name WindowIcon asset/icon.ico
${FYNE_TOOL_PATH}/bin/fyne bundle -o source/resource/Icon.go --append --name AlbumTabIcon asset/album_tab.png
${FYNE_TOOL_PATH}/bin/fyne bundle -o source/resource/Icon.go --append --name AlbumAdderOnlineIcon asset/album_adder_online.png
${FYNE_TOOL_PATH}/bin/fyne bundle -o source/resource/Icon.go --append --name MusicTabIcon asset/music_tab.png
${FYNE_TOOL_PATH}/bin/fyne bundle -o source/resource/Icon.go --append --name MusicAdderOnlineIcon asset/music_adder_online.png
${FYNE_TOOL_PATH}/bin/fyne bundle -o source/resource/Icon.go --append --name DefaultIcon asset/default.png
${FYNE_TOOL_PATH}/bin/fyne bundle -o source/resource/Icon.go --append --name RandomIcon asset/random.png
${FYNE_TOOL_PATH}/bin/fyne bundle -o source/resource/Icon.go --append --name YouTubeIcon asset/youtube.png
${FYNE_TOOL_PATH}/bin/fyne bundle -o source/resource/Icon.go --append --name BiliBiliIcon asset/bilibili.png

# Font
${FYNE_TOOL_PATH}/bin/fyne bundle -o source/resource/Font.go --package resource --name RegularFont asset/regular_font.ttf
${FYNE_TOOL_PATH}/bin/fyne bundle -o source/resource/Font.go --append --name BoldFont asset/bold_font.ttf
${FYNE_TOOL_PATH}/bin/fyne bundle -o source/resource/Font.go --append --name ItalicFont asset/italic_font.ttf
${FYNE_TOOL_PATH}/bin/fyne bundle -o source/resource/Font.go --append --name BoldItalicFont asset/bold_italic_font.ttf

rm -rf out
mkdir out
#sudo apt install xorg-dev   if "fatal error: X11/Xcursor/Xcursor.h: No such file or directory", https://github.com/go-gl/glfw/issues/129
go build -o out/meowyplayer
