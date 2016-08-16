package system

import "log"

// CheckErr ...
func CheckErr(e error, fatal bool) {
	if fatal {
		if e != nil {
			log.Fatal(e)
		}
	} else {
		if e != nil {
			panic(e)
		}
	}
}
