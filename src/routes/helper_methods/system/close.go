package system

import (
	"log"
	"net/http"
	"os"
)

// Exit ...
func Exit(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Closed")
	os.Exit(0)
}
