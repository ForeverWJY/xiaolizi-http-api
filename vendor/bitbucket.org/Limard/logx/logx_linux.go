// +build linux

package logx

import (
	"os"
	"time"
)

var (
	defaultFilePerm = os.FileMode(0666)
)

func getDefaultLogPath() string {
	return `/var/log/`
}

var LogSaveTime = 6 * 24 * time.Hour

func addNewLine(s string) string {
	l := len(s)
	if l == 0 {
		return "\n"
	}
	if s[l-1] != '\n' {
		return s + "\n"
	}
	return s
}
