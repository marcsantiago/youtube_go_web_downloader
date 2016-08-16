package main

import (
	"log"
	"net/http"
	"runtime"

	"./src/routes/helper_methods/decorators"
	"./src/routes/homepage"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func main() {
	runtime.GOMAXPROCS(MaxParallelism())
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homepage.Index)
	log.Fatal(http.ListenAndServe(":8001", context.ClearHandler(decorators.Log((myRouter)))))

	// //Load home page
	// http.HandleFunc("/", index)
	//
	// //Post Request that handles installing ffmpeg on mac and windows
	// http.HandleFunc("/run_setup", setup)
	//
	// // validate mp3, video paths
	// http.HandleFunc("/validate_mp3_path", validateMp3)
	// http.HandleFunc("/validate_video_path", validateVideo)
	//
	// // close server
	// http.HandleFunc("/close_server", exit)
	//
	// // handle download request
	// http.HandleFunc("/download", downloader)
	//
	// //Link Static JS and CSS Files
	// path, err = os.Getwd()
	// checkErr(err, true)
	//
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
