package models

// Config ...
type Config struct {
	Setup         bool
	Mp3Path       string
	Mp3PathOkay   bool
	VideoPath     string
	VideoPathOkay bool
	ValidURL      bool
	Warning       bool
	Setup         bool
}

// DownloadMode ...
type DownloadMode struct {
	Mode string
}
