package utils

//对日志库的封装，用于后期方便更换或者拓展

import (
	"github.com/golang/glog"
	"runtime"
	"strconv"
	"strings"
)

func LogInit() {

}

func LogExit() {
	glog.Flush()
}

// 内容调用
func logFileLine(depth int) (string, int) {
	_, file, line, ok := runtime.Caller(2 + depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return file, line
}

func LogTraceD(format string, args ...interface{}) {
	file, line := logFileLine(0)
	strLine := strconv.Itoa(line)
	strfmt := "-> " + file + ":" + strLine + " " + format
	glog.Debugf(strfmt, args...)
}

func LogTraceI(format string, args ...interface{}) {
	file, line := logFileLine(0)
	strLine := strconv.Itoa(line)
	strfmt := "-> " + file + ":" + strLine + " " + format
	glog.Infof(strfmt, args...)
}

func LogTraceE(format string, args ...interface{}) {
	file, line := logFileLine(0)
	strLine := strconv.Itoa(line)
	strfmt := "-> " + file + ":" + strLine + " " + format
	glog.Errorf(strfmt, args...)
}
