package earthshaker

import (
	"os"
	"fmt"
	"time"
)

const (
	SYS = iota
	ERROR
	DEBUG
	INFO
) 

type LogSys struct {
	level int
	isopen bool
	filename string
}

var g_log LogSys

func IniLog(level int, isopen bool, filename string) bool {
	g_log.level = level
	g_log.isopen = isopen
	g_log.filename = filename + "_" + time.Now().Format("2006-01-02") + ".log"
	LOG(SYS, "IniLog in ", time.Now())
	return true
}

func OpenLog() {
	g_log.isopen = true
}

func CloseLog() {
	g_log.isopen = false
}

func gen_header(level int) string {
	time := "[" + time.Now().Format("Jan 2, 2006 at 3:04pm (MST)") + "] "
	switch level {
	case SYS:
		time += "[SYS]"
	case ERROR:
		time += "[ERROR]"
	case DEBUG:
		time += "[DEBUG]"
	case INFO:
		time += "[INFO]"
	default:
		time += "[UNKNOW]"
	}
	time += " : "
	return time
}

func LOGF(level int, format string, a ...interface{}) bool {
	if g_log.isopen && level <= g_log.level {
		str := fmt.Sprintf(format, a...)
		return write_file(gen_header(level) + str + "\r\n")
	}
	return false
}

func LOG(level int, a ...interface{}) bool {
	if g_log.isopen && level <= g_log.level {
		str := fmt.Sprint(a...)
		return write_file(gen_header(level) + str + "\r\n")
	}
	return false
}

func write_file(str string) bool {
	file, err := os.OpenFile(g_log.filename, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 666)
	if err == nil {
		file.WriteString(str)
		return true
	}
	return false
}
