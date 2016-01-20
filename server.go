package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type SetupConfig struct {
	Setup bool
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	configData := SetupConfig{}
	configData.Setup = false
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	configFile := filepath.Join(path, "/config_files/setup.json")
	if _, err := os.Stat(configFile); err == nil {
		file, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Fatal(err)
		}
		temp := SetupConfig{}
		json.Unmarshal(file, &temp)
		configData.Setup = temp.Setup
	}
	t.Execute(w, &configData)
}

func setup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var path string
		var err error
		if runtime.GOOS == "windows" {
			path, err = os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			batScript := filepath.Join(path, "/scripts/install_ffmpeg.bat")

			log.Printf("Copying windows_ffmpeg contents to c:\\FFMPEG and adding the path env c:\\FFMPEG\\bin\n")
			cmd := exec.Command("cmd", "/C", batScript)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()

		} else {
			path, err = filepath.Abs("")
			if err != nil {
				log.Fatal(err)
			}
			shellScript := filepath.Join(path, "/scripts/install_ffmpeg.sh")
			log.Printf("Please be patient installing homebrew, ffpmeg, and updating take a while\n")
			cmd := exec.Command("/bin/sh", shellScript)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}

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

		// Make sure that the page doesn't change
		//http.Redirect(w, r, "http://localhost:3000/", 301)
		// send data back to ajax post
		w.Write([]byte("ok"))
	}
}

func main() {

	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", "http://localhost:3000/").Start()
	case "windows", "darwin":
		exec.Command("open", "http://localhost:3000/").Start()
	case "unix":
		exec.Command("open", "http://localhost:3000/").Start()
	default:
		log.Fatal("unsupported platform")
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", index)
	http.HandleFunc("/run_setup", setup)

	log.Println("Listening at port 3000")
	log.Fatal(http.ListenAndServe(":3000", Log(http.DefaultServeMux)))

}
