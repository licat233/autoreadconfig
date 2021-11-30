package readfile

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

var lastModiTime time.Time
var isPanic bool

//errorHandler 错误处理回调函数类型
type errorHandler func(error)

//errHandler 错误处理
func errHandler(err error) {
	if err == nil {
		return
	}
	if lastModiTime.IsZero() {
		log.Fatal(err)
	} else {
		panic(err)
	}
}

//YamlConfig 读取yaml配置文件
//@path 文件路径string
//@out 映射结构体地址，请传入&struct{}
//errEvent 错误处理回调函数
func YamlConfig(path string, out interface{}, errEvent errorHandler) {
	defer func() {
		err := recover()
		if err != nil && !isPanic {
			isPanic = true
			errEvent(errors.New(fmt.Sprint(err)))
		} else if err == nil && isPanic {
			isPanic = false
		}
	}()
	file, err := os.Open(path)
	errHandler(err)
	fileinfo, err := file.Stat()
	errHandler(err)

	if fileinfo.ModTime() != lastModiTime {
		configByte, err := ioutil.ReadAll(file)
		errHandler(err)
		err = yaml.Unmarshal(configByte, out)
		errHandler(err)
		lastModiTime = fileinfo.ModTime().Local()
	}
}
