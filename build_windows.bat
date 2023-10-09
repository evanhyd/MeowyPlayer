rm -rf out
mkdir out
go build -ldflags -H=windowsgui -o out/meowyplayer.exe
cp -r asset out