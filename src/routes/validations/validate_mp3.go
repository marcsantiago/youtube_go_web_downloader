package validations

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"../helper_methods/system"
)

// ValidateMp3 ...
func ValidateMp3(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		mp3Path := r.FormValue("folderpath")
		mp3Path = strings.TrimSpace(mp3Path)
		if _, err := os.Stat(mp3Path); err != nil {
			w.Write([]byte("not ok"))
			masterConfig.Mp3PathOkay = false
		} else {

			path, err := os.Getwd()
			system.CheckErr(err, true)

			configFolder := filepath.Join(path, "/config_files")
			if _, err := os.Stat(configFolder); err != nil {
				//for working with binary obj
				filename := os.Args[0]
				filedirectory := filepath.Dir(filename)
				path, err = filepath.Abs(filedirectory)
				system.CheckErr(err, true)

			}
			configPath := filepath.Join(path, "/config_files/folderpaths.json")
			setupConfig := make(map[string]string)

			if _, err := os.Stat(configPath); err == nil {
				temp := PathConfig{}
				file, err := ioutil.ReadFile(configPath)
				system.CheckErr(err, true)

				json.Unmarshal(file, &temp)
				setupConfig["Mp3Path"] = mp3Path
				setupConfig["VideoPath"] = temp.VideoPath
			} else {
				setupConfig["Mp3Path"] = mp3Path
			}

			obj, err := json.Marshal(setupConfig)
			system.CheckErr(err, true)

			f, err := os.Create(configPath)
			system.CheckErr(err, true)
			defer f.Close()

			f.Write(obj)
			masterConfig.Mp3PathOkay = true
			w.Write([]byte("ok"))
		}
	}
}

// ValidateVideo ...
func ValidateVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		videoPath := r.FormValue("folderpath")
		videoPath = strings.TrimSpace(videoPath)
		if _, err := os.Stat(videoPath); err != nil {
			masterConfig.VideoPathOkay = false
			w.Write([]byte("not ok"))
		} else {

			path, err := os.Getwd()
			system.CheckErr(err, true)

			configFolder := filepath.Join(path, "/config_files")
			if _, err := os.Stat(configFolder); err != nil {
				//for working with binary obj
				filename := os.Args[0]
				filedirectory := filepath.Dir(filename)
				path, err = filepath.Abs(filedirectory)
				system.CheckErr(err, true)
			}
			configPath := filepath.Join(path, "/config_files/folderpaths.json")
			setupConfig := make(map[string]string)

			if _, err := os.Stat(configPath); err == nil {
				temp := PathConfig{}
				file, err := ioutil.ReadFile(configPath)
				system.CheckErr(err, true)

				json.Unmarshal(file, &temp)
				setupConfig["Mp3Path"] = temp.Mp3Path
				setupConfig["VideoPath"] = videoPath
			} else {
				setupConfig["VideoPath"] = videoPath
			}

			obj, err := json.Marshal(setupConfig)
			system.CheckErr(err, true)

			f, err := os.Create(configPath)
			system.CheckErr(err, true)
			defer f.Close()

			f.Write(obj)
			masterConfig.VideoPathOkay = true
			w.Write([]byte("ok"))
		}
	}
}
