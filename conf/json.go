package conf

import (
	"encoding/json"
	"farm/log"
	"io/ioutil"
)

type RunMode int //运行模式
const (
	dev     RunMode = iota //开发
	test                   //测试
	release                //发布
)

var Server struct {
	HttpServer  string
	Redis_IP    string
	Redis_Name  string
	Redis_Pwd   string
	DB_IP       string
	DB_Name     string
	DB_UserName string
	DB_Pwd      string
}

func init() {
	var mode RunMode
	mode = dev
	var file_str string
	switch mode {
	case dev:
		file_str = "conf/server_dev.json"
	case test:
		file_str = "conf/server_test.json"
	case release:
		file_str = "conf/server_release.json"
	}
	data, err := ioutil.ReadFile(file_str)
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}
}
