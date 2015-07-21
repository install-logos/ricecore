package ricecore

import (
	"os"
	"strings"
)

func exists(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func expandDir(path string) string {
	path = strings.Replace(path, "~", homeDir, 1)
	if path[len(path)-1:] != "/" {
		return path + "/"
	}
	return path
}
