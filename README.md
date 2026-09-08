[![BuyMeACoffee](https://raw.githubusercontent.com/pachadotdev/buymeacoffee-badges/main/bmc-orange.svg)](https://buymeacoffee.com/unboxthecat)
# MeowyPlayer V3
## 💻It is Modern
![Customizable Music Player](https://i.imgur.com/Nd7hmwn.png)
* 🎨 Minimalistic user interface with a modern glassy design.
* 🌐 Cross-platform support for Windows, Linux, and Mac.
* 🌍 Support for multiple languages and community-created translations.

## ✈️It is Fast
![Premium Music Resource](https://i.imgur.com/ltNPn37.png)
* 🪶 Doesn't use bulky web engines like Electron, so you can save more resources for gaming.
* 📦 Self-contained with minimal external dependencies.

## ⭐It is Premium
![Personal Backup Server](https://i.imgur.com/GACZMJs.png)
* 🔍 Grand unified music search engine that supports various platforms.
* 🎵 Premium mp3 audio quality, defaulting to 320kbps.
* 💾 Automatically backup your playlists to the cloud for free!

## Build
**Requirements:**
- Go 1.25.6+
- gcc 14.2.0+
- Run `go install fyne.io/fyne/v2/cmd/fyne@latest`
- Linux Only: `sudo apt install libxcursor-dev libxinerama-dev libxrandr-dev libxi-dev libgl-dev libxxf86vm-dev`

**Compile**
- Run `go run build.go` to compile the program with optional `-release` flag.

## Self-Hosting Servers
[MeowStore: MeowyPlayer Playlist storage server](https://github.com/evanhyd/MeowStore)  
[MeowAuth: MeowyPlayer Authentication server](https://github.com/evanhyd/MeowAuth)

## Design Documents
[MeowyPlayer Specification](https://docs.google.com/document/d/1C-rM3gNhS7BC3XFaLk-tDXIVsfijMc9mltGgSUUpcKs/edit?usp=sharing)

## Localization Credits
* **fr-FR:** Steven Gong
* **ja-JP:** Mattloulou
* **ko-KO:** TERA
* **zh-CN:** UnboxTheCat
With additional supports from LLM.

## Why MeowyPlayer
It was a tranquil summer night in 2022. I was chilling in my room doing competitive programming while listening to my friend's anime music playlist on YouTube. I tried to port the playlist to my Spotify, but to my surprise, Spotify's search engine failed to fetch the music due to a regional license restriction. Lots of anime/game OST or platform specific music were simply not available to users. I was a broke student who could barely afford his tuition (via student loans), therefore I had to spend days hacking around all kinds of sketchy websites to collect that music. It was a very frustrating and painful process, and that's when I decided to start this project.

MeowyPlayer started in July 2022 and has been re-written 5 times as of August 2026, with each iteration introducing new features, better UI design, and a more-maintainable system design.

[If you like this project, feel free to drop a tip to help me survive inflation](https://www.buymeacoffee.com/unboxthecat).
