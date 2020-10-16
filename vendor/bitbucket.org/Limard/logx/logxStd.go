package logx

import (
	"io"
)

var logxSTD = New("", "")

// Trace output a [DEBUG] trace string
func Trace() {
	logxSTD.trace()
}

// Debug output a [DEBUG] string
func Debug(v ...interface{}) {
	logxSTD.debug(v...)
}

// Debugf output a [DEBUG] string with format
func Debugf(format string, v ...interface{}) {
	logxSTD.debugf(format, v...)
}

func DebugToJson(v ...interface{}) {
	logxSTD.debugToJson(v...)
}

// Info output a [INFO ] string
func Info(v ...interface{}) {
	logxSTD.info(v...)
}

// Infof output a [INFO ] string with format
func Infof(format string, v ...interface{}) {
	logxSTD.infof(format, v...)
}

// Warn output a [WARN ] string
func Warn(v ...interface{}) {
	logxSTD.warn(v...)
}

// Warnf output a [WARN ] string with format
func Warnf(format string, v ...interface{}) {
	logxSTD.warnf(format, v...)
}

// Error output a [ERROR] string
func Error(v ...interface{}) {
	logxSTD.error(v...)
}

// Errorf output a [ERROR] string with format
func Errorf(format string, v ...interface{}) {
	logxSTD.errorf(format, v...)
}

// SetLogPath set path of output log
func SetLogPath(s string) {
	logxSTD.LogPath = s
}

func SetExeName(s string) {
	logxSTD.LogName = s
}

// SetOutputFlag set output purpose(OutputFlag_File | OutputFlag_Console | OutputFlag_DbgView)
func SetOutputFlag(flag int) {
	logxSTD.OutputFlag = flag
}

// SetOutputLevel set output level.
// OutputLevel_Debug
// OutputLevel_Info
// OutputLevel_Warn
// OutputLevel_Error
// OutputLevel_Unexpected
func SetOutputLevel(level int) {
	logxSTD.OutputLevel = level
}

// SetTimeFlag set time format(Lshortfile | Ldate | Ltime)
func SetTimeFlag(flag int) {
	logxSTD.TimeFlag = flag
}

// SetConsoleOut set a writer instead of console
func SetConsoleOut(out io.Writer) {
	logxSTD.ConsoleOutWriter = out
}
