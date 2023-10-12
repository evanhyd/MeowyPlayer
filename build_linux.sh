rm -rf out
mkdir out
#sudo apt install xorg-dev   if "fatal error: X11/Xcursor/Xcursor.h: No such file or directory", https://github.com/go-gl/glfw/issues/129
go build -o out/meowyplayer
cp -r asset out
