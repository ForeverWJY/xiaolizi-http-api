package main

import (
	"encoding/json"
)

func getOSCHINA() *OSCHINAResponse {
	key := []byte("oschinaCache")
	//先查询缓存
	value, err := cache.Get(key)
	if err != nil {
		msg := getOSCHINANow()
		if msg != nil {
			bytes, _ := json.Marshal(msg)
			cache.Set(key, bytes, 60*60*2)
		}
		return msg
	} else {
		var msg = new(OSCHINAResponse)
		_ = json.Unmarshal(value, &msg)
		return msg
	}
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
