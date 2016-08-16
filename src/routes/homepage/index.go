package homepage

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"../helper_methods/system"
)

// Index ...
func Index(w http.ResponseWriter, r *http.Request) {
	// for working within go run
	path, err := os.Getwd()
	system.CheckErr(err, true)

	templatePath := filepath.Join(path, "/templates/index.html")
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		//for working with binary obj
		filename := os.Args[0]
		filedirectory := filepath.Dir(filename)
		path, err = filepath.Abs(filedirectory)
		system.CheckErr(err, true)

		templatePath := filepath.Join(path, "/templates/index.html")
		t, _ = template.ParseFiles(templatePath)

		system.CheckErr(err, true)

	}

	masterConfig.Setup = false

	configFile := filepath.Join(path, "/config_files/setup.json")
	if _, err := os.Stat(configFile); err == nil {
		file, err := ioutil.ReadFile(configFile)
		system.CheckErr(err, true)

		temp := SetupConfig{}
		json.Unmarshal(file, &temp)
		masterConfig.Setup = temp.Setup
	}

	configFile = filepath.Join(path, "/config_files/folderpaths.json")
	if _, err := os.Stat(configFile); err == nil {
		file, err := ioutil.ReadFile(configFile)
		system.CheckErr(err, true)

		temp := PathConfig{}
		json.Unmarshal(file, &temp)
		masterConfig.Mp3Path = temp.Mp3Path
		masterConfig.VideoPath = temp.VideoPath
	}
	t.Execute(w, &masterConfig)
}
