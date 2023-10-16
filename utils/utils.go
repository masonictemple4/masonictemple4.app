package utils

import (
	"encoding/json"
	"path/filepath"
	"strings"
)

func Convert[S any, T any](source S, target T) error {
	tmp, err := json.Marshal(source)
	if err != nil {
		return err
	}
	err = json.Unmarshal(tmp, &target)
	if err != nil {
		return err
	}
	return nil
}

func GrabFileName(path string) string {
	return filepath.Base(path)
}

// There might be a better way of doing this in the future. I have tried with the bytes
// using http.DetectContentType(data) and not as much help as it should be. Will have to
// research later to see if there is another way of detecting file type.
func DetectFileType(fileName string) string {
	parts := strings.Split(fileName, ".")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-1]
}
