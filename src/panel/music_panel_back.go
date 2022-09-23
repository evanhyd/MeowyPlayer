package panel

import (
	"bufio"
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"meowyplayer.com/src/custom_canvas"
	"meowyplayer.com/src/resource"
)

//Ignore cases, check for substring
func SatisfyMusicInfo(query string, data *custom_canvas.MusicInfo) bool {
	return strings.Contains(strings.ToLower(data.Title), strings.ToLower(query))
}

//Add music to the music config
func addMusicConfig(musicTitle, albumTitle string) error {

	//load music config
	configPath := resource.GetAlbumConfigPath(albumTitle)
	configFile, err := os.OpenFile(configPath, os.O_RDWR|os.O_APPEND, fs.ModePerm)

	if err != nil {
		return err
	}
	defer configFile.Close()

	//read music name from config file
	scanner := bufio.NewScanner(configFile)
	for scanner.Scan() {
		if scanner.Text() == musicTitle {
			return nil
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	//append music title to the config
	_, err = configFile.WriteString(musicTitle + "\n")
	if err != nil {
		return err
	}
	return nil
}

//Create a music file by copying to the music folder
func AddMusicFromLocal(musicPath, musicTitle, albumTitle string) error {

	destPath := resource.GetMusicPath(musicTitle)

	//check for duplicated music file
	if _, err := os.Stat(destPath); os.IsNotExist(err) {

		//copy album icon
		musicByte, err := os.ReadFile(musicPath)
		if err != nil {
			return err
		}

		err = os.WriteFile(destPath, musicByte, os.ModePerm)
		if err != nil {
			return err
		}
	}

	//add to the config
	err := addMusicConfig(musicTitle, albumTitle)
	return err
}

//Download a music file from the remote server
func AddMusicFromRemote(youtubeUrl, musicTitle, albumTitle string) error {

	//sanitize url strings
	if strings.ContainsAny(musicTitle, `<>:"/\|?*`) {
		return errors.New("music title contains forbidden character")
	}

	//parse video ID
	videoID := youtubeUrl[strings.LastIndex(youtubeUrl, "=")+1:]
	log.Printf("Youtube Video ID: %v\n", videoID)

	//generate server post data
	serverUrl := `https://www.y2mate.com/mates/analyze/ajax`
	serverParams := url.Values{}
	serverParams.Set("url", `https://www.youtube.com/watch?v=`+videoID)
	serverParams.Set("q_auto", "1")
	serverParams.Set("ajax", "1")
	log.Printf("POST Url: %v?%v\n", serverUrl, serverParams.Encode())

	//Server POST
	serverResp, err := http.Post(serverUrl, "application/x-www-form-urlencoded; charset=UTF-8", strings.NewReader(serverParams.Encode()))
	if err != nil {
		return err
	}
	defer serverResp.Body.Close()

	//scrape user ID
	serverBody, err := ioutil.ReadAll(serverResp.Body)
	if err != nil {
		return err
	}
	serverBodyStr := string(serverBody)
	beginIDStr := `var k__id = \"`
	endIDStr := `\"; var video_service =`
	userID := ""
	if begin := strings.Index(serverBodyStr, beginIDStr); begin != -1 {
		if end := strings.Index(serverBodyStr, endIDStr); end != -1 {
			userID = serverBodyStr[begin+len(beginIDStr) : end]
		}
	}
	if userID == "" {
		return errors.New("couldn't obtain user ID")
	}
	log.Printf("User ID: %v\n", userID)

	//generate converter post data
	converterUrl := `https://www.y2mate.com/mates/convert`
	converterParams := url.Values{}
	converterParams.Set("type", "youtube")
	converterParams.Set("_id", userID)
	converterParams.Set("v_id", videoID)
	converterParams.Set("ajax", "1")
	converterParams.Set("token", "")
	converterParams.Set("ftype", "mp3")
	converterParams.Set("fquality", "128")

	//Converter POST
	converterResp, err := http.Post(converterUrl, `application/x-www-form-urlencoded; charset=UTF-8`, strings.NewReader(converterParams.Encode()))
	if err != nil {
		return err
	}
	defer converterResp.Body.Close()

	//scrape file url
	converterBody, err := ioutil.ReadAll(converterResp.Body)
	if err != nil {
		return err
	}
	converterBodyStr := string(converterBody)
	beginConverterStr := `<a href=\"`
	endConverterStr := `\" rel=\"nofollow\"`
	fileUrl := ""
	if begin := strings.Index(converterBodyStr, beginConverterStr); begin != -1 {
		if end := strings.Index(converterBodyStr, endConverterStr); end != -1 {
			fileUrl = converterBodyStr[begin+len(beginConverterStr) : end]
		}
	}
	if fileUrl == "" {
		return errors.New("couldn't obtain file url")
	}
	fileUrl = strings.ReplaceAll(fileUrl, `\/`, `/`)

	//download the music file
	fileResp, err := http.Get(fileUrl)
	if err != nil {
		return err
	}
	defer fileResp.Body.Close()

	//save to local drive
	file, err := os.Create(resource.GetMusicPath(musicTitle))
	if err != nil {
		return err
	}
	defer file.Close()

	//copy to music folder
	bytesRead, err := io.Copy(file, fileResp.Body)
	if err != nil {
		return err
	}
	log.Printf("Downloaded Size: %.2v megabytes\n", float64(bytesRead)/1024.0/1024.0)

	//add to the config
	err = addMusicConfig(musicTitle, albumTitle)
	return err
}

//Remove music from an album
func RemoveMusic(albumTitle string, index int) error {

	//load music config
	configPath := resource.GetAlbumConfigPath(albumTitle)
	configFile, err := os.OpenFile(configPath, os.O_RDWR, fs.ModePerm)

	if err != nil {
		return err
	}
	defer configFile.Close()

	//read music name from config file
	musicList := make([]string, 0)
	scanner := bufio.NewScanner(configFile)
	for scanner.Scan() {
		musicList = append(musicList, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	sort.Slice(musicList, func(i, j int) bool {
		return strings.ToLower(musicList[i]) < strings.ToLower(musicList[j])
	})

	//clear the old config file
	if err := configFile.Truncate(0); err != nil {
		return err
	}
	if _, err := configFile.Seek(0, 0); err != nil {
		return err
	}

	//override the config file
	for i := range musicList {
		if i != index {
			if _, err := configFile.WriteString(musicList[i] + "\n"); err != nil {
				return err
			}
		}
	}
	return nil
}

// //Download a music file from the remote server
// func DownloadMusicFile(youtubeUrl, title string) error {

// 	//sanitize url strings
// 	if strings.ContainsAny(title, `<>:"/\|?*`) {
// 		return errors.New("music title contains forbidden character")
// 	}

// 	//parse video ID
// 	CLIPZAG := `https://clipzag.com/`
// 	videoID := youtubeUrl[strings.LastIndex(youtubeUrl, "=")+1:]
// 	log.Printf("Youtube Video ID: %v\n", videoID)

// 	//get clipzag page response
// 	clipzagResp, err := http.Get(CLIPZAG + `watch?v=` + videoID)
// 	if err != nil {
// 		return err
// 	}
// 	defer clipzagResp.Body.Close()

// 	//parse MP3 page url
// 	MP3PageUrl := ""
// 	MP3PageUrlPrefix := `<a class="btn btn-danger margin-top-10" id="mp3button" style="font-weight:bold;padding:8px 30px;" data-mp3processed target="_blank" href="`
// 	MP3PageUrlSuffix := `"`
// 	for scanner := bufio.NewScanner(clipzagResp.Body); scanner.Scan(); {
// 		line := scanner.Text()
// 		if begin := strings.Index(line, MP3PageUrlPrefix); begin != -1 {
// 			if end := strings.LastIndex(line, MP3PageUrlSuffix); end != -1 {
// 				MP3PageUrl = line[begin+len(MP3PageUrlPrefix) : end]
// 				break
// 			}
// 		}
// 	}
// 	if MP3PageUrl == "" {
// 		return errors.New("Couldn't find the MP3 page url")
// 	}

// 	//get MP3 page response
// 	MP3PageResp, err := http.Get(CLIPZAG + MP3PageUrl)
// 	if err != nil {
// 		return err
// 	}
// 	defer clipzagResp.Body.Close()

// 	//parse MP3 file query
// 	MP3FileQuery := ""
// 	MP3FileQueryPrefix := `data:{`
// 	MP3FileQuerySuffix := `'},`
// 	for scanner := bufio.NewScanner(MP3PageResp.Body); scanner.Scan(); {
// 		line := scanner.Text()
// 		if begin := strings.Index(line, MP3FileQueryPrefix); begin != -1 {
// 			if end := strings.LastIndex(line, MP3FileQuerySuffix); end != -1 {
// 				MP3FileQuery = line[begin+len(MP3FileQueryPrefix) : end]
// 				break
// 			}
// 		}
// 	}
// 	if MP3FileQuery == "" {
// 		return errors.New("Couldn't extract the MP3 file query")
// 	}
// 	MP3FileQuery = strings.ReplaceAll(MP3FileQuery, `: '`, `=`)
// 	MP3FileQuery = strings.ReplaceAll(MP3FileQuery, `', `, `&`)

// 	//get MP3 file url response
// 	MP3ENGINE := `https://mp3.clipzag.com/mp3engine?`
// 	MP3FileQueryResp, err := http.Get(MP3ENGINE + MP3FileQuery)
// 	if err != nil {
// 		return err
// 	}
// 	defer MP3FileQueryResp.Body.Close()
// 	log.Println(MP3FileQueryResp.Status)

// 	//parse MP3 file url
// 	buffer, err := io.ReadAll(MP3FileQueryResp.Body)
// 	if err != nil {
// 		return err
// 	}
// 	MP3FileUrlStr := string(buffer)
// 	MP3FileUrlPrefix := `{"error":false,"url":"\/\/`
// 	MP3FileUrlSuffix := `","exectime"`
// 	if begin := strings.Index(MP3FileUrlStr, MP3FileUrlPrefix); begin != -1 {
// 		if end := strings.Index(MP3FileUrlStr, MP3FileUrlSuffix); end != -1 {
// 			MP3FileUrlStr = `http://` + MP3FileUrlStr[begin+len(MP3FileUrlPrefix):end]
// 			MP3FileUrlStr = strings.ReplaceAll(MP3FileUrlStr, `\/`, `/`)
// 		}
// 	}
// 	log.Println(MP3FileUrlStr)

// 	//get MP3 file
// 	fileRsp, err := http.Get(MP3FileUrlStr)
// 	if err != nil {
// 		return err
// 	}
// 	defer fileRsp.Body.Close()

// 	//save to local file
// 	file, err := os.Create(filepath.Join(resource.MUSIC_PATH, title))
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	if bytesRead, err := io.Copy(file, fileRsp.Body); err == nil {
// 		log.Printf("File Size: %.2v megabytes\n", float64(bytesRead)/1024.0/1024.0)
// 		return nil
// 	} else {
// 		return err
// 	}
// }
