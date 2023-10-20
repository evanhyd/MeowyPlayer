fyne bundle -o source/resource/Icon.go --package resource --name MissingIcon asset/missing_asset.png
fyne bundle -o source/resource/Icon.go --append --name WindowIcon asset/icon.ico
fyne bundle -o source/resource/Icon.go --append --name AlbumTabIcon asset/album_tab.png
fyne bundle -o source/resource/Icon.go --append --name AlbumAdderOnlineIcon asset/album_adder_online.png
fyne bundle -o source/resource/Icon.go --append --name MusicTabIcon asset/music_tab.png
fyne bundle -o source/resource/Icon.go --append --name MusicAdderOnlineIcon asset/music_adder_online.png
fyne bundle -o source/resource/Icon.go --append --name DefaultIcon asset/default.png
fyne bundle -o source/resource/Icon.go --append --name RandomIcon asset/random.png
fyne bundle -o source/resource/Icon.go --append --name YouTubeIcon asset/youtube.png
fyne bundle -o source/resource/Icon.go --append --name BiliBiliIcon asset/bilibili.png

fyne bundle -o source/resource/Font.go --package resource --name RegularFont asset/regular_font.ttf
fyne bundle -o source/resource/Font.go --append --name BoldFont asset/bold_font.ttf
fyne bundle -o source/resource/Font.go --append --name ItalicFont asset/italic_font.ttf
fyne bundle -o source/resource/Font.go --append --name BoldItalicFont asset/bold_italic_font.ttf

rmdir /s /q out
mkdir out
go build -ldflags -H=windowsgui -o out/meowyplayer.exe
