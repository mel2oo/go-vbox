package vbox

import "strings"

// example:
// input  -> name="win7_64_2"
// output -> name, win7_64_2
func GetLineKeyVal(s string) (key string, value string) {
	list := strings.Split(s, "=")
	if len(list) != 2 {
		return
	}

	return strings.Trim(list[0], "\""), strings.Trim(list[1], "\"")
}
