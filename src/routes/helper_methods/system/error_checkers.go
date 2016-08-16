package system

import "log"

func checkErr(e error, fatal bool) {
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
