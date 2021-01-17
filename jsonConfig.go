package main

import (
	"encoding/json"
	"io/ioutil"
)

type JsonStruct struct {
}

type Config struct {
	IP         string `json:"ip"`
	Port       int    `json:"port"`
	HttpPort   int    `json:"httpPort"`
	WsPort     int    `json:"wsPort"`
	ManagerQQ  []int  `json:"managerQQ"`
	ReportTime struct {
		Qq    []int `json:"qq"`
		Group []int `json:"group"`
	} `json:"reportTime"`
	Help string `json:"help"`
}

func (jst *JsonStruct) Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logApi.Error(err.Error())
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		logApi.Error(err.Error())
		return
	}
}

func (jst *JsonStruct) Write(filename string, v interface{}) {
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		logApi.Error(err.Error())
		return
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		logApi.Error(err.Error())
		return
	}
}

func loadJsonConfig() Config {
	JsonParse := new(JsonStruct)
	v := Config{}
	JsonParse.Load("./config.json", &v)
	return v
}

func saveJsonConfig(config Config) {
	JsonParse := new(JsonStruct)
	JsonParse.Write("./config.json", config)
}