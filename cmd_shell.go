package utils

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	// "runtime"
	"strings"
	// "syscall"
	"context"
	"io/ioutil"
	"time"
)

//指定超时时间，调用外部可执行程序运行
func CmdRunWithTimeout(ctx context.Context, timeout time.Duration, strCmd string, strArgs ...string) (error, bool) {
	//添加调试手段函数
	defer PanicHandler()

	if len(strCmd) == 0 { //避免空串
		LogTraceE("the cmd exe is empty")
		return errors.New("the cmd exe empty"), false
	}
	var outputBuf bytes.Buffer
	var cmdArgs []string
	var redirectFile string //重定向文件名
	var hasRedirect bool = false
	//添加一个用于重定向输出的文件命令结构解析
	for _, strArg := range strArgs {
		if strArg == ">" || strArg == ">>" { //重定向符
			hasRedirect = true
			LogTraceI("has redirect %s", strArg)
			continue
		} else {
			if hasRedirect {
				redirectFile = strArg
				LogTraceI("the redirect file %s", strArg)
			} else {
				cmdArgs = append(cmdArgs, strArg)
			}
		}
	}
	//创建一个channel
	done := make(chan error)
	//构造
	cmd := exec.Command(strCmd, cmdArgs...)
	//有重定向文件
	if hasRedirect {
		cmd.Stdout = &outputBuf
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// if runtime.GOOS == "windows" {
	// 	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// }

	//启动执行
	err := cmd.Start()
	if err != nil {
		LogTraceE("start run:%s failed,error:%s", cmd.Path, err.Error())
		//将应用程序的标准输出和错误输出中的信息打印出来
		// LogTraceE("stdout:%s", string(stdout.Bytes()))
		// LogTraceE("stderr:%s", string(stderr.Bytes()))
		return err, false
	}
	go func() {
		done <- cmd.Wait()
	}()

	// 创建一个超时计时器
	chTimeout := time.After(timeout)

	select {
	case <-ctx.Done():
		if err = cmd.Process.Kill(); err != nil {
			LogTraceE("failed to kill:%s,error:%s", cmd.Path, err.Error())
		}
		go func() {
			//防止因为kill以后,wait goroutine里面的done因为没有接受，导致无法写数据到done中
			<-done //allow wait goroutine to exit
		}()
		LogTraceI("process:%s killed,because upper exit", cmd.Path)
		return errors.New("process run but upper exit"), true
	case <-chTimeout:
		if err = cmd.Process.Kill(); err != nil {
			LogTraceE("failed to kill:%s,error:%s", cmd.Path, err.Error())
		}
		go func() {
			//防止因为kill以后,wait goroutine里面的done因为没有接受，导致无法写数据到done中
			<-done //allow wait goroutine to exit
		}()
		LogTraceI("process:%s killed,because timeout", cmd.Path)
		return errors.New("process run timeout"), true
	case err = <-done:
		if err != nil {
			return err, false
		}

	}
	//如果有重定向文件，需要把输出输出到重定向文件中去
	if hasRedirect {
		if err = ioutil.WriteFile(redirectFile, outputBuf.Bytes(), 0777); err != nil {
			LogTraceE("write redirect file :%s failed", redirectFile)
			return err, false
		}
		LogTraceI("write redirect file :%s success", redirectFile)
	}
	return nil, false
}

//指定超时时间，调用外部可执行程序运行
func EasyCmdRunWithTimeout(ctx context.Context, timeout time.Duration, strCmd string, strArgs string) (error, bool) {
	// 对参数进行切片
	Args := strings.Split(strArgs, " ")
	LogTraceI("len(args)=%d;%+v", len(Args), Args)
	var correctArgs []string
	startIndex := -1
	joinFlag := false
	for index, argv := range Args {
		//引号
		if strings.Contains(argv, "\"") {
			if strings.Index(argv, "\"") == strings.LastIndex(argv, "\"") { //仅仅含单个的时候
				if joinFlag {
					joinArg := strings.Join(Args[startIndex:index+1], " ")
					LogTraceI("joinArgs:%s", joinArg)
					joinArg = strings.Replace(joinArg, "\"", "", -1)
					LogTraceI("replace joinArgs:%s", joinArg)
					correctArgs = append(correctArgs, joinArg)
				}
				startIndex = index
				joinFlag = !joinFlag
			} else if argv != " " && !joinFlag {
				argv = strings.Replace(argv, "\"", "", -1)
				correctArgs = append(correctArgs, argv)
			}
		} else if argv != " " && !joinFlag {
			correctArgs = append(correctArgs, argv)
		}
	}
	LogTraceI("len(correctArgs)=%d;%+v", len(correctArgs), correctArgs)
	// 执行合成操作
	err, isTimeout := CmdRunWithTimeout(ctx, timeout, strCmd, correctArgs...)
	if err != nil {
		return err, isTimeout
	}
	return err, isTimeout
}

//运行外部程序,通过给chCmd发送消息来实现结束这个外部进程
func CmdRunProcess(chKill chan string, strCmd string, strArgs ...string) error {
	//添加调试手段函数
	defer PanicHandler()

	if len(strCmd) == 0 { //避免空串
		LogTraceE("the cmd exe is empty")
		return errors.New("the cmd exe empty")
	}
	//创建一个channel
	done := make(chan error)
	//构造
	cmd := exec.Command(strCmd, strArgs...)
	// var stdout, stderr bytes.Buffer
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// if runtime.GOOS == "windows" {
	// 	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// }

	//启动执行
	err := cmd.Start()
	if err != nil {
		LogTraceE("start run:%s failed,error:%s\n", cmd.Path, err.Error())
		//将应用程序的标准输出和错误输出中的信息打印出来
		// LogTraceE("stdout:%s", string(stdout.Bytes()))
		// LogTraceE("stderr:%s", string(stderr.Bytes()))
		return err
	}
	//起一个线程来等待进程执行完成后的返回值
	go func() {
		done <- cmd.Wait()
	}()

	LogTraceI("the start process :%s; pid:%d", cmd.Path, cmd.Process.Pid)
	select {
	case strChCmd := <-chKill:
		LogTraceI("receive %s form chKill", strChCmd)
		if err = cmd.Process.Kill(); err != nil {
			LogTraceE("failed to kill:%s,error:%s\n", cmd.Path, err.Error())
		}
		select {
		case exitErr := <-done: //等kill 以后，进程成功退出
			if cmd.ProcessState.Exited() {
				LogTraceI("the process exited ok\n")
			}
			LogTraceI("the process exited err:%s\n", exitErr.Error())
		}
		LogTraceI("process:%s killed\n", cmd.Path)
		return nil
	case err = <-done:
		LogTraceE("process:%s err:%s\n", err.Error())
		return err
	}
	return nil
}

//测试
func TestCmdRun() {
	defer PanicHandler()
	// // aerender
	// strCmd := "E:\\Program Files\\Adobe\\Adobe After Effects CC 2017\\Support Files\\aerender.exe"
	// // strArgs := " -project \"F:\\test_file\\AE\\012\\012.aepx\" -comp \"WRITE_USR_VIDEO_PRODUCT_001\" -s 1 -e 85 -output \"D:\\001.mov\" "
	// strArgs := []string{"-project", "\"F:\\test_file\\AE\\012\\012.aepx\"", "-comp", "\"WRITE_USR_VIDEO_PRODUCT_001\"", "-s", "1", "-e", "85", "-output", "\"D:\\001.mov\""}

	// ffmpeg
	strCmd := "F:\\test_file\\ffmpeg.exe "
	// strArgs := " -i F:\\test_file\\ff.mp4 -ignore_loop 0 -i F:\\test_file\\gif_4.gif -filter_complex [0:v][1:v]overlay=shortest=1 -vcodec libx264 -an -y F:\\test_file\\ff_gif_4_b.mp4"
	strArgs := []string{"-i", "F:\\test_file\\ff.mp4", "-ignore_loop", "0", "-i", "F:\\test_file\\gif_4.gif",
		"-filter_complex", "[0:v][1:v]overlay=shortest=1", "-vcodec", "libx264", "-an", "-y", "F:\\test_file\\ff_gif_4_b.mp4"}

	timeout := 5 * time.Minute
	err, isTimeout := CmdRunWithTimeout(context.Background(), timeout, strCmd, strArgs...)
	// err, isTimeout := CmdRunWithTimeout(timeout, strCmd, strArgs)
	if err != nil {
		LogTraceI("%s,%d", err.Error(), isTimeout)
	} else {
		LogTraceI("cmd run ok")
	}
}
