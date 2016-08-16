package downloader

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"../helper_methods/system"
)

// Downloader ...
func Downloader(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var wg sync.WaitGroup

		urlData := r.FormValue("url")
		// singleMode := r.FormValue("SingleMode")
		mp3Mode := r.FormValue("MP3Mode")

		// clean up the data
		urlSplit := strings.Split(urlData, "\n")
		for i := 0; i < len(urlSplit); i++ {
			urlSplit[i] = strings.TrimSpace(urlSplit[i])
		}
		log.Println("Downloading data")
		for _, url := range urlSplit {
			wg.Add(1)
			ExcecuteDownload(url, mp3Mode, &wg)
		}
		wg.Wait()

		//checking for vidoes and moving
		// path = strings.Replace(path, "desktop", "Desktop", -1)
		// log.Println(path)
		log.Println("Moving videos or mp3s to the current folder")
		videos := system.CheckExt(".m4a")
		videos = append(videos, system.CheckExt(".webm")...)
		videos = append(videos, system.CheckExt(".mp4")...)
		videos = append(videos, system.CheckExt(".3gp")...)
		videos = append(videos, system.CheckExt(".flv")...)
		for _, vid := range videos {
			currentVideoPath := filepath.Join(path, vid)
			newPath := filepath.Join(masterConfig.VideoPath, vid)
			newPath = strings.Replace(newPath, "#", "", -1)
			err := os.Rename(currentVideoPath, newPath)
			system.CheckErr(err, true)
		}

		//moving mp3s
		mp3s := system.CheckExt(".mp3")
		for _, m := range mp3s {
			currentMp3Path := filepath.Join(path, m)
			// newPath := filepath.Join(masterConfig.Mp3Path, m)
			newPath = strings.Replace(newPath, "#", "", -1)
			err := os.Rename(currentMp3Path, newPath)
			checkErr(err, true)
		}
		log.Println("Complete Please Check Your Folders")
		w.Write([]byte("ok"))
		return
	}
}
