package main

import (
	"encoding/json"
	"net/url"
	"strings"
)

func getWeather(city string) string {
	key := []byte("weatherCache" + city)
	//先查询缓存
	value, err := cache.Get(key)
	if err != nil {
		logApi.Error(err.Error())
		weather := getWeatherNow(city)
		if weather != "" {
			_ = cache.Set(key, []byte(weather + "\ncache for 2 hours"), 60*60*2)
		}
		return weather
	} else {
		return string(value[:])
	}
}

func getWeatherNow(city string) string {
	v := url.Values{}
	v.Add("city", city)
	body := v.Encode()
	bytes := get("http://wthrcdn.etouch.cn/weather_mini?" + body)
	if bytes == nil {
		logApi.Error("获取天气出错")
		return ""
	}
	var resp = new(WeatherResponse)
	_ = json.Unmarshal(bytes, resp)

	if resp.Desc == "OK" && resp.Status == 1000 {
		var str []string
		str = append(str, "城市："+resp.Data.City+"，温度："+resp.Data.Wendu+"℃\n")
		for _, v := range resp.Data.Forecast {
			index := strings.Index(v.Fengli, "]")
			str = append(str, v.Date+"："+v.Type+" "+v.Low+" "+v.High+" "+v.Fengxiang+v.Fengli[9:index]+"\n")
		}
		str = append(str, "感冒情况："+resp.Data.Ganmao)
		logApi.Debug(strings.Join(str, ""))
		return strings.Join(str, "")
	}
	return ""
}

type WeatherResponse struct {
	Data struct {
		Yesterday struct {
			Date string `json:"date"`
			High string `json:"high"`
			Fx   string `json:"fx"`
			Low  string `json:"low"`
			Fl   string `json:"fl"`
			Type string `json:"type"`
		} `json:"yesterday"`
		City     string `json:"city"`
		Forecast []struct {
			Date      string `json:"date"`
			High      string `json:"high"`
			Fengli    string `json:"fengli"`
			Low       string `json:"low"`
			Fengxiang string `json:"fengxiang"`
			Type      string `json:"type"`
		} `json:"forecast"`
		Ganmao string `json:"ganmao"`
		Wendu  string `json:"wendu"`
	} `json:"data"`
	Status int    `json:"status"`
	Desc   string `json:"desc"`
}
