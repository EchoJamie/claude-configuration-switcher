package util

import (
	"os"
	"path/filepath"
	"strings"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ExpandPath 展开路径中的 ~ 为用户主目录
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			home = os.Getenv("HOME")
		}
		if home != "" {
			path = filepath.Join(home, path[1:])
		}
	}
	return path
}
