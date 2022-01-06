package utils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

func IsExistFileInfo(name string) error {
	//获取文件信息
	fileInfo, err := os.Stat(name)
	if err != nil {
		LogTraceE("the "+name+" not exist,err is ", err.Error())
		return err
	}

	if fileInfo.Size() > 0 || fileInfo.IsDir() { //这个文件大小>0或者它是一个目录文件，表示这个文件存在
		return nil
	}

	return errors.New("file not exist or invalid")
}

//根据当前日期来创建文件夹
func CreateDateDir(basePath string) string {
	foldName := time.Now().Format("20180913")
	foldPath := basePath + foldName + "/"
	if _, err := os.Stat(foldPath); os.IsNotExist(err) {
		//分成两步，先创建文件夹，然后再修改权限
		os.Mkdir(foldPath, 0777)
		os.Chmod(foldPath, 0777)
	}
	return foldPath
}

//获取文件的md5值
func GetFileMd5(path string) string {
	var md5V string
	f, err := os.Open(path)
	if err != nil {
		LogTraceE("Open %s", err.Error())
		return md5V
	}

	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		LogTraceE("Copy %s", err.Error())
		return md5V
	}

	md5V = fmt.Sprintf("%x", md5hash.Sum(nil))

	return md5V
}

func GetFileSize(fileName string) (int64, error) {
	fileInfo, err := os.Stat(fileName)
	if nil != err {
		LogTraceE("the %s %s", fileName, err.Error())
		return 0, err
	}

	return fileInfo.Size(), nil
}

//文件拷贝
func CopyFile(destPath, srcPath string) error {
	if destPath == srcPath {
		return nil
	}
	os.Remove(destPath)

	//取目录
	dirs := path.Dir(destPath)
	if "" == dirs {
		LogTraceE("the destPath dir is empty")
		return errors.New("the destPath dir is empty")
	}

	//创建目录
	if err := os.MkdirAll(dirs, 0755); nil != err {
		LogTraceE("MkdirAll error：%s,%s", dirs, err.Error())
		return err
	}
	//
	src, err := os.Open(srcPath)
	if nil != err {
		return err
	}
	defer src.Close()

	n, err := GetFileSize(srcPath)
	if nil != err {
		return err
	}

	dst, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE, 0644)
	if nil != err {
		return err
	}
	defer dst.Close()

	written, err := io.CopyN(dst, src, n)
	if nil != err {
		return err
	}

	if n != written {
		LogTraceE("copy size error, src[%d], dst[%d].", n, written)
		return errors.New("copy size error")
	}

	return nil
}

func ExpiredDirClean(dir string, expiredTime time.Duration) error {
	dir_info, err := os.Stat(dir)
	if nil != err {
		LogTraceE("dir[%s] stat:%s", dir, err.Error())
		return err
	}
	dir_mod_time := dir_info.ModTime() //返回目录最后修改的时间
	now := time.Now()
	if now.Sub(dir_mod_time) > expiredTime { //目录时间过期了，就清理掉
		if err := os.RemoveAll(dir); err != nil {
			LogTraceE("dir[%s] remove:%s", dir, err.Error())
			return err
		}
	}
	return nil
}
