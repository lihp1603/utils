package utils

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

var (
	//定义panic调试输出的方式
	panicToStd = flag.Bool("panic_trace", false, "panic info to standard error instead of file")
)

// PanicTrace trace panic stack info.
func PanicTrace() string {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, 2048)
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return string(stack)
}

func PanicTraceEx() {
	if err := recover(); err != nil {
		log.Println(PanicTrace())
	}
}

func PanicDump() {
	if err := recover(); err != nil {
		exeName := os.Args[0] //获取程序名称

		now := time.Now()  //获取当前时间
		pid := os.Getpid() //获取进程ID

		time_str := now.Format("2006-01-02.15-04-05")                     //设定时间格式
		fname := fmt.Sprintf("%s-%d-%s-dump.log", exeName, pid, time_str) //保存错误信息文件名:程序名-进程ID-当前时间（年月日时分秒）
		fmt.Println("dump to file ", fname)

		f, err := os.Create(fname)
		if err != nil {
			return
		}
		defer f.Close()
		f.WriteString(fmt.Sprintf("%v\r\n", err)) //输出panic信息
		f.WriteString("========\r\n")
		f.WriteString(string(debug.Stack())) //输出堆栈信息
	}
}

func PanicHandler() {
	if *panicToStd {
		PanicTraceEx()
	} else {
		PanicDump()
	}
}
