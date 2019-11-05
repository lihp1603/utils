package utils

//对日志库的封装，用于后期方便更换或者拓展

import (
	"github.com/golang/glog"
	"github.com/lihp1603/scribe"
	"runtime"
	"strconv"
	"strings"
)

var (
	gScribeTrace *scribe.ScribeTrace = nil
)

type LogLevel int32

const (
	DebugLog LogLevel = iota
	InfoLog
	ErrorLog
)

var loglevelName = []string{
	DebugLog: "DEBUG",
	InfoLog:  "INFO",
	ErrorLog: "ERROR"}

func LogInit(log_scribe interface{}) {
	if log_scribe != nil {
		if scribe_trace, ok := log_scribe.(*scribe.ScribeTrace); ok {
			gScribeTrace = scribe_trace
		}
	}
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
	if gScribeTrace != nil {
		gScribeTrace.WriteScribeEx(file, strLine, "DEBUG", format, args...)
	}
}

func LogTraceI(format string, args ...interface{}) {
	file, line := logFileLine(0)
	strLine := strconv.Itoa(line)
	strfmt := "-> " + file + ":" + strLine + " " + format
	glog.Infof(strfmt, args...)
	if gScribeTrace != nil {
		gScribeTrace.WriteScribeEx(file, strLine, "INFO", format, args...)
	}
}

func LogTraceE(format string, args ...interface{}) {
	file, line := logFileLine(0)
	strLine := strconv.Itoa(line)
	strfmt := "-> " + file + ":" + strLine + " " + format
	glog.Errorf(strfmt, args...)
	if gScribeTrace != nil {
		gScribeTrace.WriteScribeEx(file, strLine, "ERROR", format, args...)
	}
}

//添加一个综合的函数
func LogTrace(log_level LogLevel, log_scribe interface{}, format string, args ...interface{}) {
	var scribeTrace *scribe.ScribeTrace = nil
	if log_scribe != nil {
		if scribe_trace, ok := log_scribe.(*scribe.ScribeTrace); ok {
			scribeTrace = scribe_trace
		}
	}
	//获取打印位置代码的文件和行数
	file, line := logFileLine(0)
	strLine := strconv.Itoa(line)
	strfmt := "-> " + file + ":" + strLine + " " + format
	//日志级别判断
	switch log_level {
	case DebugLog:
		glog.Debugf(strfmt, args...)
	case InfoLog:
		glog.Infof(strfmt, args...)
	case ErrorLog:
		glog.Errorf(strfmt, args...)
	}

	if scribeTrace != nil {
		scribeTrace.WriteScribeEx(file, strLine, loglevelName[log_level], format, args...)
	}
	return
}
