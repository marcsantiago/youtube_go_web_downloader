package main

import (
	"log"
	"net/http"
	"runtime"

	"./src/routes/downloader"
	"./src/routes/helper_methods/decorators"
	"./src/routes/helper_methods/system"
	"./src/routes/homepage"
	"./src/routes/setup"
	"./src/routes/validations"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func main() {
	runtime.GOMAXPROCS(system.MaxParallelism())
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homepage.Index)
	//Post Request that handles installing ffmpeg on mac and windows
	myRouter.HandleFunc("/run_setup", setup.Setup)
	// validate mp3, video paths
	myRouter.HandleFunc("/validate_mp3_path", validations.ValidateMp3)
	myRouter.HandleFunc("/validate_video_path", validations.ValidateVideo)
	// handle download request
	myRouter.HandleFunc("/download", downloader.Downloader)
	// serve static contents
	myRouter.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))
	log.Fatal(http.ListenAndServe(":8001", context.ClearHandler(decorators.Log((myRouter)))))

	// jsPath := filepath.Join(path, "static/js")
	// if _, err := os.Stat(jsPath); err != nil {
	// 	//for working with binary obj
	// 	filename := os.Args[0]
	// 	filedirectory := filepath.Dir(filename)
	// 	path, err = filepath.Abs(filedirectory)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	jsPath = filepath.Join(path, "static/js")
	// }
	//
	// macPath = filepath.Join(path, "mac/")
	// windowsPath = filepath.Join(path, "windows_binaries/")
	//
	// js := http.FileServer(http.Dir(jsPath))
	// http.Handle("/static/js/", http.StripPrefix("/static/js/", js))
	//
	// cssPath := filepath.Join(path, "static/css")
	// css := http.FileServer(http.Dir(cssPath))
	// http.Handle("/static/css/", http.StripPrefix("/static/css/", css))
	//
	// fontPath := filepath.Join(path, "static/fonts")
	// font := http.FileServer(http.Dir(fontPath))
	// http.Handle("/static/fonts/", http.StripPrefix("/static/fonts/", font))
	//
	// log.Println("Listening at port 3000")
	// log.Fatal(http.ListenAndServe(":3000", Log(http.DefaultServeMux)))

}
