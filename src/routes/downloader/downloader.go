package downloader

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Downloader ...
func Downloader(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		urlData := r.FormValue("url")
		//singleMode := r.FormValue("SingleMode")
		mp3Mode := r.FormValue("MP3Mode")

		//clean up the urls
		log.Println("Checking Urls")
		urlSplit := strings.Split(urlData, "\n")
		for i := 0; i < len(urlSplit); i++ {
			urlSplit[i] = strings.TrimSpace(urlSplit[i])
		}
		// check to make sure they are valid youtube urls
		masterConfig.ValidUrl = true
		for _, url := range urlSplit {
			if checkUrl(url) == false {
				masterConfig.ValidUrl = false
				w.Write([]byte("not ok"))
				return
			}
		}

		log.Println("Downloading data")
		for _, url := range urlSplit {
			wg.Add(1)
			downloaderfile(url, mp3Mode)
		}
		wg.Wait()

		//checking for vidoes and moving
		path = strings.Replace(path, "desktop", "Desktop", -1)
		log.Println(path)
		log.Println("Moving videos or mp3s to the current folder")
		videos := checkExt(".m4a")
		videos = append(videos, checkExt(".webm")...)
		videos = append(videos, checkExt(".mp4")...)
		videos = append(videos, checkExt(".3gp")...)
		videos = append(videos, checkExt(".flv")...)
		for _, vid := range videos {
			currentVideoPath := filepath.Join(path, vid)
			newPath := filepath.Join(masterConfig.VideoPath, vid)
			newPath = strings.Replace(newPath, "#", "", -1)
			err = os.Rename(currentVideoPath, newPath)
			checkErr(err, true)
		}

		//moving mp3s
		mp3s := checkExt(".mp3")
		for _, m := range mp3s {
			currentMp3Path := filepath.Join(path, m)
			newPath := filepath.Join(masterConfig.Mp3Path, m)
			newPath = strings.Replace(newPath, "#", "", -1)
			err = os.Rename(currentMp3Path, newPath)
			checkErr(err, true)
		}
		log.Println("Complete Please Check Your Folders")
		w.Write([]byte("ok"))
	}
}
