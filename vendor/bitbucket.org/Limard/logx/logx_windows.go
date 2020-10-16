// +build windows

package logx

import (
	"os"
	"path/filepath"
	"time"
)

var (
	defaultFilePerm = os.FileMode(0666)
)

func getDefaultLogPath() string {
	return filepath.Dir(os.Args[0])
}

// LogSaveTime ...
var LogSaveTime = 6 * 24 * time.Hour

func addNewLine(s string) string {
	l := len(s)
	if l == 0 {
		return "\r\n"
	}
	if l == 1 {
		return s + "\r\n"
	}
	if s[l-2] == '\r' && s[l-1] == '\n' {
		return s
	}
	if s[l-1] == '\r' || s[l-1] == '\n' {
		return s[:l-1] + "\r\n"
	}
	return s + "\r\n"
}
