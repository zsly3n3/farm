package conf

import (
	"encoding/json"
	"farm/log"
	"io/ioutil"
)

type RunMode int //运行模式
const (
	Dev     RunMode = iota //开发
	Test                   //测试
	Release                //发布
)

const Mode = Dev

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

	var file_str string
	switch Mode {
	case Dev:
		file_str = "conf/server_dev.json"
	case Test:
		file_str = "conf/server_test.json"
	case Release:
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
