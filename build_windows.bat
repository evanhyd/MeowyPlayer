rmdir /s /q out
mkdir out
go build -ldflags -H=windowsgui -o out/meowyplayer.exe
xcopy /E /I /Y asset out\asset