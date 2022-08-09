package util

import (
	"path"
	"path/filepath"
	"runtime"
)

func String(value string, defaultValue string) string {
	if len(value) > 0 {
		return value
	}
	return defaultValue
}

func Root() string {
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		return filepath.Dir(path.Join(path.Dir(filename), "../../"))
	}
	return ""
}
