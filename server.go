package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func setup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		if runtime.GOOS == "windows" {
			path, err := os.Getwd()
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
			path, err := filepath.Abs("")
			if err != nil {
				log.Fatal(err)
			}
			shellScript := filepath.Join(path, "/scripts/install_ffmpeg.sh")
			log.Printf("Please be patient installing homebrew, ffpmeg, and updating take a while")
			cmd := exec.Command("/bin/sh", shellScript)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
		// Make sure that the page doesn't change
		//http.Redirect(w, r, "http://localhost:3000/", 301)
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

	assets := http.StripPrefix("/", http.FileServer(http.Dir("static/")))
	http.Handle("/", assets)
	http.HandleFunc("/run_setup", setup)

	log.Println("Listening at port 3000")
	log.Fatal(http.ListenAndServe(":3000", Log(http.DefaultServeMux)))

}
