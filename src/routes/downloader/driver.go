package downloader

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// ExcecuteDownload ...
func ExcecuteDownload(url string, mp3Mode string) {
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
