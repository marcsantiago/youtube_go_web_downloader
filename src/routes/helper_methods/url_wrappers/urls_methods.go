package urls

import "strings"

// CheckURL ...
func CheckURL(url string) bool {
	//Modify this method to except more type of links
	if strings.Contains(url, "https://") == true {
		return true
	}
	return false
}
