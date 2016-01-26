package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"path/filepath"
	"runtime"
)

//////////////////////
////DATA STRUCTURES
//////////////////////
type Config struct {
	Setup         bool
	Mp3Path       string
	Mp3PathOkay   bool
	VideoPath     string
	VideoPathOkay bool
	ValidUrl      bool
	Warning       bool
}

type DownloaderInfo struct {
	SingleMode bool
	MP3Mode    bool
}

type SetupConfig struct {
	Setup bool
}

type PathConfig struct {
	Mp3Path   string
	VideoPath string
}

//////////////////////
////GLOBAL VARIABLES
//////////////////////
var (
	masterConfig = Config{}
	macPath      string
	windowsPath  string
	platform     string
	path         string
	err          error
	wg           sync.WaitGroup
)

//////////////////////
////HELPER METHODS
//////////////////////
func checkExt(ext string) []string {
	pathS, err := os.Getwd()
	checkErr(err, false)

	var files []string
	filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r && (strings.Contains(f.Name(), "pyc") == false || strings.Contains(f.Name(), "mp3.py") == false) {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	return files
}

func checkErr(e error, fatal bool) {
	if fatal {
		if e != nil {
			log.Fatal(e)
		}
	} else {
		if e != nil {
			panic(e)
		}
	}
}

func checkUrl(url string) bool {
	//Modify this method to except more type of links
	if strings.Contains(url, "https://") == true {
		return true
	}
	return false
}

func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

//////////////////////
////MAIN ROUTE (HOMEPAGE)
//////////////////////
func index(w http.ResponseWriter, r *http.Request) {
	// for working within go run
	path, err := os.Getwd()
	checkErr(err, true)

	templatePath := filepath.Join(path, "/templates/index.html")
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		//for working with binary obj
		filename := os.Args[0]
		filedirectory := filepath.Dir(filename)
		path, err = filepath.Abs(filedirectory)
		checkErr(err, true)

		templatePath := filepath.Join(path, "/templates/index.html")
		t, _ = template.ParseFiles(templatePath)

		checkErr(err, true)

	}

	masterConfig.Setup = false

	configFile := filepath.Join(path, "/config_files/setup.json")
	if _, err := os.Stat(configFile); err == nil {
		file, err := ioutil.ReadFile(configFile)
		checkErr(err, true)

		temp := SetupConfig{}
		json.Unmarshal(file, &temp)
		masterConfig.Setup = temp.Setup
	}

	configFile = filepath.Join(path, "/config_files/folderpaths.json")
	if _, err := os.Stat(configFile); err == nil {
		file, err := ioutil.ReadFile(configFile)
		checkErr(err, true)

		temp := PathConfig{}
		json.Unmarshal(file, &temp)
		masterConfig.Mp3Path = temp.Mp3Path
		masterConfig.VideoPath = temp.VideoPath
	}
	t.Execute(w, &masterConfig)
}

//////////////////////
////CLOSE SERVER
//////////////////////
func exit(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Closed")
	os.Exit(0)
}

//////////////////////
////SETUP FFMPEG
//////////////////////
func setup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		path, err := os.Getwd()
		checkErr(err, true)

		batScript := filepath.Join(path, "/scripts/install_ffmpeg.bat")
		if _, err := os.Stat(batScript); err != nil {
			//for working with binary obj
			filename := os.Args[0]
			filedirectory := filepath.Dir(filename)
			path, err = filepath.Abs(filedirectory)
			checkErr(err, true)
			batScript = filepath.Join(path, "/scripts/install_ffmpeg.bat")
		}

		switch runtime.GOOS {
		case "darwin", "unix":
			shellScript := filepath.Join(path, "/scripts/install_ffmpeg.sh")
			log.Printf("Please be patient installing homebrew, ffpmeg, youtube-dl, and updating take a while\n")
			cmd := exec.Command("/bin/sh", shellScript)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		case "windows":
			log.Printf("Copying windows_binaries contents to c:\\FFMPEG and adding the path env c:\\FFMPEG\\bin\n")
			cmd := exec.Command("cmd", "/C", batScript)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		default:
			log.Fatal("unsupported platform")
			os.Exit(1)
		}

		configPath := filepath.Join(path, "/config_files/setup.json")
		setupConfig := make(map[string]bool)
		setupConfig["Setup"] = true

		obj, err := json.Marshal(setupConfig)
		checkErr(err, true)

		f, err := os.Create(configPath)
		checkErr(err, true)
		defer f.Close()

		f.Write(obj)
		w.Write([]byte("ok"))
	}
}

//////////////////////
////ENSURE MP3 DOWNLOAD PATH
//////////////////////
func validateMp3(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		mp3Path := r.FormValue("folderpath")
		mp3Path = strings.TrimSpace(mp3Path)
		if _, err := os.Stat(mp3Path); err != nil {
			w.Write([]byte("not ok"))
			masterConfig.Mp3PathOkay = false
		} else {

			path, err := os.Getwd()
			checkErr(err, true)

			configFolder := filepath.Join(path, "/config_files")
			if _, err := os.Stat(configFolder); err != nil {
				//for working with binary obj
				filename := os.Args[0]
				filedirectory := filepath.Dir(filename)
				path, err = filepath.Abs(filedirectory)
				checkErr(err, true)

			}
			configPath := filepath.Join(path, "/config_files/folderpaths.json")
			setupConfig := make(map[string]string)

			if _, err := os.Stat(configPath); err == nil {
				temp := PathConfig{}
				file, err := ioutil.ReadFile(configPath)
				checkErr(err, true)

				json.Unmarshal(file, &temp)
				setupConfig["Mp3Path"] = mp3Path
				setupConfig["VideoPath"] = temp.VideoPath
			} else {
				setupConfig["Mp3Path"] = mp3Path
			}

			obj, err := json.Marshal(setupConfig)
			checkErr(err, true)

			f, err := os.Create(configPath)
			checkErr(err, true)
			defer f.Close()

			f.Write(obj)
			masterConfig.Mp3PathOkay = true
			w.Write([]byte("ok"))
		}
	}
}

//////////////////////
////ENSURE VIDEO DOWNLOAD PATH
//////////////////////
func validateVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		videoPath := r.FormValue("folderpath")
		videoPath = strings.TrimSpace(videoPath)
		if _, err := os.Stat(videoPath); err != nil {
			masterConfig.VideoPathOkay = false
			w.Write([]byte("not ok"))
		} else {

			path, err := os.Getwd()
			checkErr(err, true)

			configFolder := filepath.Join(path, "/config_files")
			if _, err := os.Stat(configFolder); err != nil {
				//for working with binary obj
				filename := os.Args[0]
				filedirectory := filepath.Dir(filename)
				path, err = filepath.Abs(filedirectory)
				checkErr(err, true)
			}
			configPath := filepath.Join(path, "/config_files/folderpaths.json")
			setupConfig := make(map[string]string)

			if _, err := os.Stat(configPath); err == nil {
				temp := PathConfig{}
				file, err := ioutil.ReadFile(configPath)
				checkErr(err, true)

				json.Unmarshal(file, &temp)
				setupConfig["Mp3Path"] = temp.Mp3Path
				setupConfig["VideoPath"] = videoPath
			} else {
				setupConfig["VideoPath"] = videoPath
			}

			obj, err := json.Marshal(setupConfig)
			checkErr(err, true)

			f, err := os.Create(configPath)
			checkErr(err, true)
			defer f.Close()

			f.Write(obj)
			masterConfig.VideoPathOkay = true
			w.Write([]byte("ok"))
		}
	}
}

//////////////////////
////DOWNLOAD DRIVER
//////////////////////
func downloader(w http.ResponseWriter, r *http.Request) {
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

//////////////////////
////DOWNLOAD CONTENT
//////////////////////
func downloaderfile(url string, mp3Mode string) {
	if mp3Mode == "true" {
		if platform == "unix" {
			log.Printf("Downloading mp3 %s\n", url)
			tool := fmt.Sprintf("youtube-dl --extract-audio --audio-format mp3 -o \"%%(title)s.%%(ext)s \" " + url)
			cmd := exec.Command("/bin/sh", "-c", tool)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		} else {
			tool := fmt.Sprintf("youtube-dl.exe --ignore-errors --extract-audio --audio-format mp3 -o \"%%(title)s.%%(ext)s \" " + url)
			cmd := exec.Command("cmd", "/C", tool)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
	} else {
		if platform == "unix" {
			log.Printf("Downloading video %s\n", url)
			cmd := exec.Command("/bin/sh", "-c", "youtube-dl -f 22 "+url)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		} else {
			cmd := exec.Command("cmd", "/C", "youtube-dl.exe -f 22 "+url)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
	}
	wg.Done()
}

//////////////////////
////MAIN (WHERE ROUTES ARE LOCATED)
//////////////////////
func main() {
	//set number of cores to use to max
	runtime.GOMAXPROCS(MaxParallelism())
	masterConfig.ValidUrl = true
	masterConfig.Mp3PathOkay = true
	masterConfig.VideoPathOkay = true
	masterConfig.Warning = false

	switch runtime.GOOS {
	case "darwin", "unix":
		platform = "unix"
		exec.Command("open", "http://localhost:3000/").Start()
	case "windows":
		platform = "windows"
		//chrome
		cmd := exec.Command("C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe", "http://localhost:3000/")
		err := cmd.Run()
		if err != nil {
			//firefox
			cmd = exec.Command("C:\\Program Files (x86)\\Mozilla Firefox\\firefox.exe", "http://localhost:3000/")
			err = cmd.Run()
			//internet explore
			if err != nil {
				exec.Command("C:\\Program Files\\Internet Explorer\\iexplore.exe", "http://localhost:3000/").Start()
				masterConfig.Warning = true
			}
		}
	default:
		log.Fatal("unsupported platform")
		os.Exit(1)
	}

	//Load home page
	http.HandleFunc("/", index)

	//Post Request that handles installing ffmpeg on mac and windows
	http.HandleFunc("/run_setup", setup)

	// validate mp3, video paths
	http.HandleFunc("/validate_mp3_path", validateMp3)
	http.HandleFunc("/validate_video_path", validateVideo)

	// close server
	http.HandleFunc("/close_server", exit)

	// handle download request
	http.HandleFunc("/download", downloader)

	//Link Static JS and CSS Files
	path, err = os.Getwd()
	checkErr(err, true)

	jsPath := filepath.Join(path, "static/js")
	if _, err := os.Stat(jsPath); err != nil {
		//for working with binary obj
		filename := os.Args[0]
		filedirectory := filepath.Dir(filename)
		path, err = filepath.Abs(filedirectory)
		if err != nil {
			log.Fatal(err)
		}
		jsPath = filepath.Join(path, "static/js")
	}

	macPath = filepath.Join(path, "mac/")
	windowsPath = filepath.Join(path, "windows_binaries/")

	js := http.FileServer(http.Dir(jsPath))
	http.Handle("/static/js/", http.StripPrefix("/static/js/", js))

	cssPath := filepath.Join(path, "static/css")
	css := http.FileServer(http.Dir(cssPath))
	http.Handle("/static/css/", http.StripPrefix("/static/css/", css))

	fontPath := filepath.Join(path, "static/fonts")
	font := http.FileServer(http.Dir(fontPath))
	http.Handle("/static/fonts/", http.StripPrefix("/static/fonts/", font))

	log.Println("Listening at port 3000")
	log.Fatal(http.ListenAndServe(":3000", Log(http.DefaultServeMux)))

}
