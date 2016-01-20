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

	configData := SetupConfig{}
	configData.Setup = false

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

	//Load home page
	http.HandleFunc("/", index)

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

	//Post Request that handles installing ffmpeg on mac and windows
	http.HandleFunc("/run_setup", setup)

	log.Println("Listening at port 3000")
	log.Fatal(http.ListenAndServe(":3000", Log(http.DefaultServeMux)))

}
