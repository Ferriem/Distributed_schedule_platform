package utils

import "os"

// Exists checks whether the file/folder in the given path exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		// Check if the error indicates that the file or directory doesn't exist
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// Ext returns the file name extension used by path.
// The extension is the suffix beginning at the final dot
// in the final slash-separated element of path;
// it is empty if there is no dot.
func Ext(path string) string {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
	}
	return ""
}
