package main

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io/ioutil"
	"net/http"
	"strings"
)

func getOSCHINA() []string {
	key := []byte("oschinaCache")
	//先查询缓存
	value, err := cache.Get(key)
	if err != nil {
		//msg := getOSCHINANow()
		msg := getOSChinaResult("industry")
		if msg != nil {
			msg = append(msg, "cache for 2 hours")
			strs := strings.Join(msg, "\n")
			cache.Set(key, []byte(strs + "\n"), 60*60*2)
			return msg
		}
		return nil
	} else {
		s := string(value[:])
		return strings.Split(s, "\n")
	}
}

//func getOSChinaResult(newType string) {
//	url := "https://www.oschina.net/action/ajax/get_more_news_list?p=1&newsType=" + newType
//	doc, err := htmlquery.LoadURL(url)
//	if err != nil {
//		logApi.Error("获取OSCHINA出错")
//		//return nil
//	}
//	list := htmlquery.Find(doc, "//a")
//	logApi.Debugf("属性值 %v", list)
//	for _, n := range list {
//		fmt.Println(htmlquery.SelectAttr(n, "href")) // output @href value
//	}
//}

func getURL(url string) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")

	resp, err := (&http.Client{}).Do(req)
	//resp, err := http.Get(serviceUrl + "/topic/query/false/lsj")
	if err != nil {
		logApi.Error(err.Error())
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

func getOSChinaResult(newType string) []string {
	url := "https://www.oschina.net/action/ajax/get_more_news_list?p=1&newsType=" + newType
	logApi.Debug("请求地址：", url)
	bytes := getURL(url)
	if bytes == nil {
		logApi.Error("获取OSCHINA出错")
		//return nil
	}
	//logApi.Debug(fmt.Sprintf("%s", bytes))
	doc, err := htmlquery.Parse(strings.NewReader(fmt.Sprintf("%s", bytes)))
	if err != nil {
		logApi.Error(err.Error())
	}
	list := htmlquery.Find(doc, "//div[@class='main-info box-aw']/a")
	var str []string
	for _, n := range list {
		title := htmlquery.InnerText(n)
		href := htmlquery.SelectAttr(n, "href")
		if strings.Index(href, "/") == 0 {
			href = "https://www.oschina.net" + href
		}
		str = append(str, title+"->"+href)
	}
	logApi.Debug(len(str))
	return str
	//return strings.Join(str, "\n")
}

func getOSCHINANow() *OSCHINAResponse {
	bytes := get("http://blog.wjyup.com/test/api/oschina/str/")
	if bytes == nil {
		logApi.Error("获取OSCHINA出错")
		return nil
	}
	var resp = new(OSCHINAResponse)
	_ = json.Unmarshal(bytes, resp)

	//for _, v := range *resp {
	//	logApi.Debug(v)
	//}
	return resp
}

type OSCHINAResponse []string
