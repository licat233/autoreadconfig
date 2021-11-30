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

//文件读取结构体
//@LastModiTime 最后修改时间
//@IsPanic 是否Panic
type ReadFile struct {
	LastModiTime time.Time
	IsPanic      bool
}

//错误处理回调函数类型
type ErrorHandler func(error)

//errHandler 错误处理
func (r *ReadFile) errHandler(err error) {
	if err == nil {
		return
	}
	if r.LastModiTime.IsZero() {
		log.Fatal(err)
	} else {
		panic(err)
	}
}

//YamlConfig 读取yaml配置文件
//@path 文件路径string
//@out 映射结构体地址，请传入&struct{}
//errEvent 错误处理回调函数
func (r *ReadFile) YamlConfig(path string, out interface{}, errEvent ErrorHandler) {
	defer func() {
		err := recover()
		if err != nil && !r.IsPanic {
			r.IsPanic = true
			errEvent(errors.New(fmt.Sprint(err)))
		} else if err == nil && r.IsPanic {
			r.IsPanic = false
		}
	}()
	file, err := os.Open(path)
	r.errHandler(err)
	fileinfo, err := file.Stat()
	r.errHandler(err)

	if fileinfo.ModTime() != r.LastModiTime {
		configByte, err := ioutil.ReadAll(file)
		r.errHandler(err)
		err = yaml.Unmarshal(configByte, out)
		r.errHandler(err)
		r.LastModiTime = fileinfo.ModTime().Local()
	}
}
