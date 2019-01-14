package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// 加载配置文件
func Load(name string, value interface{}) (err error) {
	var f *os.File
	if f, err = os.Open(name); err != nil {
		LogTraceI("open conf file failed,err is %s", err.Error())
		return
	}

	defer f.Close()

	var data []byte

	if data, err = ioutil.ReadAll(f); err != nil {
		LogTraceI("conf file load filed,err is %s", err.Error())
		return
	}

	if err = json.Unmarshal(data, value); err != nil {
		LogTraceI("conf file json fmt decoding failed,err is %s", err.Error())
		return
	}

	return
}
