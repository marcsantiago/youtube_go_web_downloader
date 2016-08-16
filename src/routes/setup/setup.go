package setups

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Setup ...
func Setup(w http.ResponseWriter, r *http.Request) {
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
