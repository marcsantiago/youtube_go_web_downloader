package main

import (
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
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

	log.Println("Listening at port 3000")
	log.Fatal(http.ListenAndServe(":3000", Log(http.DefaultServeMux)))

}
