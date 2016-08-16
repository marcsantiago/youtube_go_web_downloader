package system

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CheckExt ...
func CheckExt(ext string) []string {
	paths, err := os.Getwd()
	checkErr(err, false)

	var files []string
	filepath.Walk(paths, func(path string, f os.FileInfo, _ error) error {
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
