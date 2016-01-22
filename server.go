package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"path/filepath"
	"runtime"
)

type SetupConfig struct {
	Setup bool
}

type PathConfig struct {
	Mp3Path   string
	VideoPath string
}

type Config struct {
	Setup     bool
	Mp3Path   string
	VideoPath string
	ValidUrl  bool
}

var masterConfig = Config{}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	// for working within go run
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	templatePath := filepath.Join(path, "/templates/index.html")
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		//for working with binary obj
		filename := os.Args[0]
		filedirectory := filepath.Dir(filename)
		path, err = filepath.Abs(filedirectory)
		if err != nil {
			log.Fatal(err)
		}
		templatePath := filepath.Join(path, "/templates/index.html")
		t, _ = template.ParseFiles(templatePath)
	}

	masterConfig.Setup = false

	configFile := filepath.Join(path, "/config_files/setup.json")
	if _, err := os.Stat(configFile); err == nil {
		file, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Fatal(err)
		}

		temp := SetupConfig{}
		json.Unmarshal(file, &temp)
		masterConfig.Setup = temp.Setup
	}

	configFile = filepath.Join(path, "/config_files/folderpaths.json")
	if _, err := os.Stat(configFile); err == nil {
		file, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Fatal(err)
		}

		temp := PathConfig{}
		json.Unmarshal(file, &temp)
		masterConfig.Mp3Path = temp.Mp3Path
		masterConfig.VideoPath = temp.VideoPath
	}
	t.Execute(w, &masterConfig)
}

func exit(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Closed")
	os.Exit(0)
}

func setup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		path, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		batScript := filepath.Join(path, "/scripts/install_ffmpeg.bat")
		if _, err := os.Stat(batScript); err != nil {
			//for working with binary obj
			filename := os.Args[0]
			filedirectory := filepath.Dir(filename)
			path, err = filepath.Abs(filedirectory)
			if err != nil {
				log.Fatal(err)
			}
			batScript = filepath.Join(path, "static/js")
		}

		log.Printf("Copying windows_ffmpeg contents to c:\\FFMPEG and adding the path env c:\\FFMPEG\\bin\n")
		cmd := exec.Command("cmd", "/C", batScript)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

		shellScript := filepath.Join(path, "/scripts/install_ffmpeg.sh")
		log.Printf("Please be patient installing homebrew, ffpmeg, and updating take a while\n")
		cmd = exec.Command("/bin/sh", shellScript)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

		configPath := filepath.Join(path, "/config_files/setup.json")
		setupConfig := make(map[string]bool)
		setupConfig["Setup"] = true

		obj, err := json.Marshal(setupConfig)
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(configPath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		f.Write(obj)

		w.Write([]byte("ok"))
	}
}

func validateMp3(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		mp3Path := r.FormValue("folderpath")
		mp3Path = strings.TrimSpace(mp3Path)
		if _, err := os.Stat(mp3Path); err != nil {
			w.Write([]byte("not ok"))
		} else {

			path, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			configFolder := filepath.Join(path, "/config_files")
			if _, err := os.Stat(configFolder); err != nil {
				//for working with binary obj
				filename := os.Args[0]
				filedirectory := filepath.Dir(filename)
				path, err = filepath.Abs(filedirectory)
				if err != nil {
					log.Fatal(err)
				}
			}
			configPath := filepath.Join(path, "/config_files/folderpaths.json")
			setupConfig := make(map[string]string)

			if _, err := os.Stat(configPath); err == nil {
				temp := PathConfig{}
				file, err := ioutil.ReadFile(configPath)
				if err != nil {
					log.Fatal(err)
				}
				json.Unmarshal(file, &temp)
				setupConfig["Mp3Path"] = mp3Path
				setupConfig["VideoPath"] = temp.VideoPath
			} else {
				setupConfig["Mp3Path"] = mp3Path
			}

			obj, err := json.Marshal(setupConfig)
			if err != nil {
				log.Fatal(err)
			}

			f, err := os.Create(configPath)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			f.Write(obj)
			w.Write([]byte("ok"))
		}
	}
}

func validateVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		videoPath := r.FormValue("folderpath")
		videoPath = strings.TrimSpace(videoPath)
		if _, err := os.Stat(videoPath); err != nil {
			w.Write([]byte("not ok"))
		} else {

			path, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			configFolder := filepath.Join(path, "/config_files")
			if _, err := os.Stat(configFolder); err != nil {
				//for working with binary obj
				filename := os.Args[0]
				filedirectory := filepath.Dir(filename)
				path, err = filepath.Abs(filedirectory)
				if err != nil {
					log.Fatal(err)
				}
			}
			configPath := filepath.Join(path, "/config_files/folderpaths.json")
			setupConfig := make(map[string]string)

			if _, err := os.Stat(configPath); err == nil {
				temp := PathConfig{}
				file, err := ioutil.ReadFile(configPath)
				if err != nil {
					log.Fatal(err)
				}
				json.Unmarshal(file, &temp)
				setupConfig["Mp3Path"] = temp.Mp3Path
				setupConfig["VideoPath"] = videoPath
			} else {
				setupConfig["VideoPath"] = videoPath
			}

			obj, err := json.Marshal(setupConfig)
			if err != nil {
				log.Fatal(err)
			}

			f, err := os.Create(configPath)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			f.Write(obj)
			w.Write([]byte("ok"))
		}
	}
}

type DownloaderInfo struct {
	SingleMode bool
	MP3Mode    bool
}

func checkUrl(url string) bool {
	if strings.Contains(url, "https://www.youtube.com/watch") == true || strings.Contains(url, "https://www.youtube.com/playlist") == true {
		return true
	}
	return false
}

func downloader(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		urlData := r.FormValue("url")
		//singleMode := r.FormValue("SingleMode")
		//mp3Mode := r.FormValue("MP3Mode")

		//clean up the urls
		urlSplit := strings.Split(urlData, "\n")
		for i := 0; i < len(urlSplit); i++ {
			urlSplit[i] = strings.TrimSpace(urlSplit[i])
		}

		masterConfig.ValidUrl = true
		for _, url := range urlSplit {
			if checkUrl(url) == false {
				masterConfig.ValidUrl = false
				w.Write([]byte("not ok"))
			}
		}

	}
}

func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func main() {
	//set number of cores to use to max
	runtime.GOMAXPROCS(MaxParallelism())
	masterConfig.ValidUrl = true

	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", "http://localhost:3000/").Start()
	case "darwin":
		exec.Command("open", "http://localhost:3000/").Start()
	case "windows":
		exec.Command("C:\\Program Files\\Internet Explorer\\iexplore.exe", "http://localhost:3000/").Start()
	case "unix":
		exec.Command("open", "http://localhost:3000/").Start()
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
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

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
