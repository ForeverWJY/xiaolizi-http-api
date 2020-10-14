package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func getLoginQQ() int {
	qq := get(fmt.Sprintf("http://%v/getlogonqq", *addr))
	if qq == nil {
		log.Println("获取QQ号出错")
		return 0
	}
	type loginQQRet struct {
		Ret string `json:"ret"`
	}
	var ret = new(loginQQRet)
	_ = json.Unmarshal(qq, &ret)
	fmt.Println(ret.Ret)

	var event map[string]interface{}
	_ = json.Unmarshal([]byte(ret.Ret), &event)
	/*使用键输出地图值 */
	for qq := range event["QQlist"].(map[string]interface{}) {
		fmt.Println(qq)
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
	log.Printf("send %v\n", bodystr)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(bodystr))
	if err != nil {
		// handle error
		log.Printf("POST请求出错了 %v", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		log.Printf("ioutil.ReadAll出错了 %v", err.Error())
	}
	//fmt.Println(string(body))
	log.Printf("post返回值:%s\n", body)
	return body
}

func get(url string) []byte {
	log.Printf("get url %v\n", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("get error %v", err.Error())
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("get返回值:%s\n", body)
	return body
}