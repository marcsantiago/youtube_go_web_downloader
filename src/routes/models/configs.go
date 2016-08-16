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
}

// DownloaderInfo ...
type DownloaderInfo struct {
	SingleMode bool
	MP3Mode    bool
}

// SetupConfig ...
type SetupConfig struct {
	Setup bool
}

// PathConfig ...
type PathConfig struct {
	Mp3Path   string
	VideoPath string
}
