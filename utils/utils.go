package utils

import "encoding/json"

// generate a function that inputs an interface and outputs a json
func JSONstringify(i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(b)
}
