package logx

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// const value
const (
	OutputFlag_File = 1 << iota
	OutputFlag_Console

	OutputLevel_Debug      = 100
	OutputLevel_Info       = 200
	OutputLevel_Warn       = 300
	OutputLevel_Error      = 400
	OutputLevel_Unexpected = 500

	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

//type Loggerx struct
type Loggerx struct {
	OutFile          *os.File
	LastError        error
	FilePerm         os.FileMode
	LineMaxLength    int    // 一行最大的长度
	LogPath          string // log的保存目录
	LogName          string // log的文件名，默认为程序名
	OutputFlag       int    // 输出Flag
	OutputLevel      int    // 输出级别
	TimeFlag         int    // properties
	MaxLogNumber     int    // 最多log文件个数
	ContinuousLog    bool   // 连续在上一个文件中输出，适用于经常被调用启动的程序日志
	ConsoleOutWriter io.Writer

	logCounter int
	Prefix     []byte // Prefix to write at beginning of each line
	muFile     sync.Mutex
}

func New(path, name string) *Loggerx {
	l := &Loggerx{
		FilePerm:         defaultFilePerm,
		LineMaxLength:    1024,
		LogPath:          path,
		OutputFlag:       OutputFlag_File | OutputFlag_Console,
		OutputLevel:      OutputLevel_Debug,
		TimeFlag:         Lshortfile | Ldate | Ltime,
		MaxLogNumber:     3,
		LogName:          name,
		logCounter:       0,
		ContinuousLog:    true,
		ConsoleOutWriter: os.Stdout,
	}

	if len(l.LogPath) == 0 {
		l.LogPath = getDefaultLogPath()
	}

	if len(l.LogName) == 0 {
		n, _ := exec.LookPath(os.Args[0])
		l.LogName = filepath.Base(n)
	}

	//
	type configFile struct {
		OutputLevel string
		OutputFlag  []string
	}
	buf, e := ioutil.ReadFile("log.json")
	if e == nil {
		var c1 configFile
		json.Unmarshal(buf, &c1)

		if len(c1.OutputFlag) != 0 {
			l.OutputFlag = 0
			for _, f := range c1.OutputFlag {
				switch strings.ToLower(f) {
				case "file":
					l.OutputFlag |= OutputFlag_File
				case "console":
					l.OutputFlag |= OutputFlag_Console
				}
			}
		}

		if c1.OutputLevel != "" {
			switch strings.ToLower(c1.OutputLevel) {
			case "debug", "dbg":
				l.OutputLevel = OutputLevel_Debug
			case "info":
				l.OutputLevel = OutputLevel_Info
			case "warn", "warning":
				l.OutputLevel = OutputLevel_Warn
			case "error", "err":
				l.OutputLevel = OutputLevel_Error
			case "unexpected":
				l.OutputLevel = OutputLevel_Unexpected
			}
		}
	}

	return l
}

func funcName() string {
	funcName := ""
	pc, _, _, ok := runtime.Caller(3)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
		s := strings.Split(funcName, ".")
		funcName = s[len(s)-1]
	}
	return funcName
}

func (t *Loggerx) Trace() {
	t.trace()
}

func (t *Loggerx) trace() {
	if t.OutputLevel > OutputLevel_Debug {
		return
	}

	t.output(fmt.Sprintf("[TRACE][%s]", funcName()))
}

// Debug output a [DEBUG] string
func (t *Loggerx) Debug(v ...interface{}) {
	t.debug(v...)
}

func (t *Loggerx) debug(v ...interface{}) {
	if t.OutputLevel > OutputLevel_Debug {
		return
	}
	t.output(fmt.Sprintf(`[DEBUG][%s]%s`, funcName(), fmt.Sprint(v...)))
}

// Debugf output a [DEBUG] string with format
func (t *Loggerx) Debugf(format string, v ...interface{}) {
	t.debugf(format, v...)
}

func (t *Loggerx) debugf(format string, v ...interface{}) {
	if t.OutputLevel > OutputLevel_Debug {
		return
	}
	t.output(fmt.Sprintf(fmt.Sprintf(`[DEBUG][%s]%s`, funcName(), format), v...))
}

func (t *Loggerx) DebugToJson(v ...interface{}) {
	t.debugToJson(v...)
}

func (t *Loggerx) debugToJson(v ...interface{}) {
	if t.OutputLevel > OutputLevel_Debug {
		return
	}
	ss := []string{`[DEBUG]`, `[` + funcName() + `]`}
	for _, sub := range v {
		switch sub.(type) {
		case string:
			ss = append(ss, sub.(string))
		default:
			buf, _ := json.Marshal(sub)
			ss = append(ss, string(buf))
		}
	}
	t.output(strings.Join(ss, ""))
}

// Info output a [INFO ] string
func (t *Loggerx) Info(v ...interface{}) {
	t.info(v...)
}

func (t *Loggerx) info(v ...interface{}) {
	if t.OutputLevel > OutputLevel_Info {
		return
	}
	t.output(fmt.Sprintf(`[INFO ][%s]%s`, funcName(), fmt.Sprint(v...)))
}

// Infof output a [INFO ] string with format
func (t *Loggerx) Infof(format string, v ...interface{}) {
	t.infof(format, v...)
}

func (t *Loggerx) infof(format string, v ...interface{}) {
	if t.OutputLevel > OutputLevel_Info {
		return
	}
	t.output(fmt.Sprintf(fmt.Sprintf(`[INFO ][%s]%s`, funcName(), format), v...))
}

// Warn output a [WARN ] string
func (t *Loggerx) Warn(v ...interface{}) {
	t.warn(v...)
}

func (t *Loggerx) warn(v ...interface{}) {
	if t.OutputLevel > OutputLevel_Warn {
		return
	}
	t.output(fmt.Sprintf(`[WARN ][%s]%s`, funcName(), fmt.Sprint(v...)))
}

// Warnf output a [WARN ] string with format
func (t *Loggerx) Warnf(format string, v ...interface{}) {
	t.warnf(format, v...)
}

func (t *Loggerx) warnf(format string, v ...interface{}) {
	if t.OutputLevel > OutputLevel_Warn {
		return
	}
	t.output(fmt.Sprintf(fmt.Sprintf(`[WARN ][%s]%s`, funcName(), format), v...))
}

// Error output a [ERROR] string
func (t *Loggerx) Error(v ...interface{}) {
	t.error(v...)
}

func (t *Loggerx) error(v ...interface{}) {
	if t.OutputLevel > OutputLevel_Error {
		return
	}
	t.output(fmt.Sprintf(`[ERROR][%s]%s`, funcName(), fmt.Sprint(v...)))
}

// Errorf output a [ERROR] string with format
func (t *Loggerx) Errorf(format string, v ...interface{}) {
	t.errorf(format, v...)
}

func (t *Loggerx) errorf(format string, v ...interface{}) {
	if t.OutputLevel > OutputLevel_Error {
		return
	}
	t.output(fmt.Sprintf(fmt.Sprintf(`[ERROR][%s]%s`, funcName(), format), v...))
}

func (t *Loggerx) getFileHandle() error {
	e := os.MkdirAll(t.LogPath, 0777)
	if e != nil {
		t.LastError = e
		return e
	}

	files := make([]string, 0)
	filepath.Walk(t.LogPath, func(fPath string, fInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fInfo.IsDir() || !strings.HasPrefix(fInfo.Name(), t.LogName+`.`) || !strings.HasSuffix(fInfo.Name(), ".log") {
			return nil
		}
		if time.Now().Sub(fInfo.ModTime()) > LogSaveTime {
			os.Remove(fPath)
			return nil
		}
		files = append(files, fInfo.Name())
		return nil
	})
	for _, value := range t.getNeedDeleteLogfile(files) {
		os.Remove(t.LogPath + value)
	}

	if t.ContinuousLog {
		f := t.getNewestLogfile(files)
		if len(f) > 0 {
			filename := filepath.Join(t.LogPath, f)
			fi, e := os.Stat(filename)
			if e == nil && fi.Size() < 1024*1024*3 {
				t.OutFile, e = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, t.FilePerm)
				if e != nil {
					fmt.Println("logx:", e)
				} else {
					t.OutFile.Write([]byte("\r\n==================================================\r\n"))
				}
			} else if e != nil {
				fmt.Println("logx:", e)
			}
		}
	}
	if t.OutFile == nil {
		filename := filepath.Join(t.LogPath, t.LogName+`.`+time.Now().Format(`060102_150405`)+`.log`)
		t.OutFile, e = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, t.FilePerm)
	}
	if e != nil {
		fmt.Println("logx:", e)
		t.LastError = e
		return e
	}
	return nil
}

// 获取同名Log中最老的数个
func (t *Loggerx) getNeedDeleteLogfile(filesName []string) []string {
	if len(filesName) < t.MaxLogNumber {
		return nil
	}
	sort.Strings(filesName)
	return filesName[0 : len(filesName)-t.MaxLogNumber]
}

// 获取同名Log中最新的一个
func (t *Loggerx) getNewestLogfile(filesName []string) string {
	if len(filesName) == 0 {
		return ""
	}
	sort.Strings(filesName)
	return filesName[len(filesName)-1]
}

func (t *Loggerx) renewLogFile() (e error) {
	if t.OutFile != nil && t.logCounter < 100 {
		t.logCounter++
		return nil
	}
	t.logCounter = 1

	t.muFile.Lock()
	defer t.muFile.Unlock()

	if t.OutFile == nil {
		e = t.getFileHandle()
		if e != nil {
			return e
		}
	}

	fi, _ := t.OutFile.Stat()
	if fi.Size() > 1024*1024*3 {
		t.OutFile.Close()
		t.OutFile = nil
		e = t.getFileHandle()
		if e != nil {
			return e
		}
	}

	if t.OutFile == nil {
		return fmt.Errorf("OutFile is nil")
	}
	return nil
}

func (t *Loggerx) output(s string) {
	buf := t.makeStr(4, s)

	if t.OutputFlag&OutputFlag_File != 0 {
		e := t.renewLogFile()
		if e != nil {
			es := addNewLine(e.Error())
			if t.ConsoleOutWriter != nil {
				t.ConsoleOutWriter.Write([]byte(es))
			}
			if strings.Contains(e.Error(), "permission denied") {
				t.OutputFlag &= ^OutputFlag_File
			}
		} else {
			t.muFile.Lock()
			t.OutFile.Write(buf)
			t.muFile.Unlock()
		}
	}

	if t.OutputFlag&OutputFlag_Console != 0 && t.ConsoleOutWriter != nil {
		t.ConsoleOutWriter.Write(buf)
	}
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// formatHeader writes log header to buf in following order:
//   * l.Prefix (if it's not blank),
//   * date and/or time (if corresponding flags are provided),
//   * file and line number (if corresponding flags are provided).
func (t *Loggerx) formatHeader(buf *[]byte, tm time.Time, file string, line int) {
	*buf = append(*buf, t.Prefix...)
	if t.TimeFlag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if t.TimeFlag&LUTC != 0 {
			tm = tm.UTC()
		}
		if t.TimeFlag&Ldate != 0 {
			year, month, day := tm.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if t.TimeFlag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := tm.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if t.TimeFlag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, tm.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if t.TimeFlag&(Lshortfile|Llongfile) != 0 {
		if t.TimeFlag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

func (t *Loggerx) makeStr(calldepth int, s string) []byte {
	now := time.Now() // get this early.
	var file string
	var line int
	if t.TimeFlag&(Lshortfile|Llongfile) != 0 {
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
	}
	var buf []byte
	t.formatHeader(&buf, now, file, line)

	// limit max length
	if len(s) > t.LineMaxLength {
		buf = append(buf, s[:t.LineMaxLength]...)
		buf = append(buf, []byte(" ...")...)
	} else {
		buf = append(buf, s...)
	}

	if len(s) < 2 || s[len(s)-2] != '\r' || s[len(s)-1] != '\n' {
		buf = append(buf, '\r', '\n')
	}
	return buf
}

// For log
func (t *Loggerx) Print(v ...interface{}) {
	t.debug(v...)
}
func (t *Loggerx) Println(v ...interface{}) {
	t.debug(v...)
}
func (t *Loggerx) Printf(format string, v ...interface{}) {
	t.debugf(format, v...)
}
func (t *Loggerx) Fatal(v ...interface{}) {
	t.error(v...)
}
