package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func reply(reqMsg *ReceiveMessage, msg string) {
	switch reqMsg.Type {
	case "PrivateMsg":
		sendPrivateMsg(LoginQQ, reqMsg.FromQQ.UIN, msg)
		break
	case "GroupMsg":
		sendGroupMsg(LoginQQ, reqMsg.FromGroup.GIN, msg)
		break
	}
}

type loginQQRet struct {
	Ret string `json:"ret"`
}
func getLoginQQ() int {
	qq := get(fmt.Sprintf("http://%v/getlogonqq", *addr))
	if qq == nil {
		logApi.Debug("获取QQ号出错")
		return 0
	}
	var ret = new(loginQQRet)
	_ = json.Unmarshal(qq, &ret)
	logApi.Debug(ret.Ret)

	var event map[string]interface{}
	_ = json.Unmarshal([]byte(ret.Ret), &event)
	for qq := range event["QQlist"].(map[string]interface{}) {
		logApi.Debug("当前登录QQ：", qq)
		atoi, _ := strconv.Atoi(qq)
		return atoi
	}
	return 0
}

func sendPrivateMsg(fromqq int, toqq int, text string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	r.Form.Add("text", fmt.Sprintf("%v", text))
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/sendprivatemsg", *addr), bodystr)
}

func sendGroupMsg(fromqq int, toGroup int, text string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("togroup", fmt.Sprintf("%v", toGroup))
	r.Form.Add("text", fmt.Sprintf("%v", text))
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/sendgroupmsg", *addr), bodystr)
}

func postFormData(url string, bodystr string) []byte {
	logApi.Debugf("send %v", bodystr)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(bodystr))
	if err != nil {
		logApi.Debugf("POST请求出错了 %v", err.Error())
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logApi.Debugf("ioutil.ReadAll出错了 %v", err.Error())
		return nil
	}
	_ = resp.Body.Close()
	logApi.Debugf("post返回值:%s", body)
	return body
}

func get(url string) []byte {
	logApi.Debugf("get url %v", url)
	resp, err := http.Get(url)
	if err != nil {
		logApi.Debugf("get error %v", err.Error())
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logApi.Debugf("get read error %v", err.Error())
		return nil
	}
	_ = resp.Body.Close()
	logApi.Debugf("get返回值:%s", body)
	return body
}