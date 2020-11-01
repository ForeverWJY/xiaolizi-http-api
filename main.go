package main

import (
	"bitbucket.org/Limard/logx"
	"encoding/json"
	"fmt"
	"github.com/coocood/freecache"
	"github.com/gorilla/websocket"
	"github.com/robfig/cron"
	"io/ioutil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"
)

//var addr = flag.String("addr", "192.168.58.134:10429", "http service address")

//CRON Expression Format
//A cron expression represents a set of times, using 5 space-separated fields.
//Field name   | Mandatory? | Allowed values  | Allowed special characters
//----------   | ---------- | --------------  | --------------------------
//Minutes      | Yes        | 0-59            | * / , -
//Hours        | Yes        | 0-23            | * / , -
//Day of month | Yes        | 1-31            | * / , - ?
//Month        | Yes        | 1-12 or JAN-DEC | * / , -
//Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?

var (
	LoginQQ int //当前登录QQ
	logApi  = logx.New(".", "xiaolizi-http")
	//pongTime = 5 * time.Second
	//pingTime = 5 * time.Second
	cacheSize   = 100 * 1024 * 1024 // In bytes, where 1024 * 1024 represents a single Megabyte, and 100 * 1024*1024 represents 100 Megabytes.
	cache       = freecache.NewCache(cacheSize)
	addr        string
	managerQQ   = make(map[int]int)
	autoTimeMap = make(map[int]*cron.Cron)
)

type JsonStruct struct {
}

type Config struct {
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	ManagerQQ []int  `json:"managerQQ"`
	Help      string `json:"help"`
}

func (jst *JsonStruct) Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}

func loadManager(qqList []int) {
	for _, i := range qqList {
		fmt.Println(i)
		managerQQ[i] = i
	}
}

func isMangerQQ(rm *ReceiveMessage, f func()) {
	i := managerQQ[rm.FromQQ.UIN]
	if i > 0 {
		if f != nil {
			f()
		}
	} else {
		reply(rm, "仅管理员可操作")
	}
}

func main() {
	//flag.Parse()
	//log.SetFlags(0)

	JsonParse := new(JsonStruct)
	v := Config{}
	JsonParse.Load("./config.json", &v)

	addr = fmt.Sprintf("%v:%v", v.IP, v.Port)

	if v.ManagerQQ != nil {
		loadManager(v.ManagerQQ)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	logApi.Debugf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logApi.Error("dial:", err)
	}

	c.SetPingHandler(func(appData string) error {
		_ = c.WriteMessage(websocket.PongMessage, []byte{})
		return nil
	})

	c.SetPongHandler(func(appData string) error {
		_ = c.WriteMessage(websocket.PingMessage, []byte{})
		return nil
	})

	//c.SetReadDeadline(time.Now().Add(pongTime))

	defer c.Close()

	//获取当前QQ
	LoginQQ = getLoginQQ()

	testFunc()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logApi.Debug("read:", err)
				return
			}
			go func() {
				logApi.Debugf("recv: %s", message)
				//if len(message) > 0 {
				//	for i, ch := range message {
				//		switch {
				//		case ch > '~':
				//			message[i] = ' '
				//		case ch == '\r':
				//		case ch == '\n':
				//		case ch == '\t':
				//		case ch < ' ':
				//			message[i] = ' '
				//		}
				//	}
				//}
				//logApi.Debugf("recv1: %s", message)
				//msg := fmt.Sprintf("%s", message)
				//logApi.Debug(msg)
				var rm = new(ReceiveMessage)
				err = json.Unmarshal(message, &rm)
				if err != nil {
					logApi.Debug("转实体类出错", err.Error())
					return
				}
				//干其他事情
				logApi.Debugf("type: %v", rm.Type)
				//忽略自己的消息
				if rm.FromQQ.UIN == LoginQQ {
					return
				}
				groupQQ := rm.FromGroup.GIN
				if rm.Msg.Text == "help" {
					reply(rm, v.Help)
				} else if rm.Msg.Text == "开启报时" {
					isMangerQQ(rm, func() {
						newCron := cron.New()
						_, err2 := newCron.AddFunc("0 * * * *", func() {
							reply(rm, time.Now().Format("2006-01-02 15:04:05"))
						})
						if err2 != nil {
							logx.Error(err2.Error())
						}
						newCron.Start()
						autoTimeMap[groupQQ] = newCron
						reply(rm, "本群自动报时已开启")
					})
				} else if rm.Msg.Text == "关闭报时" {
					isMangerQQ(rm, func() {
						delCron := autoTimeMap[groupQQ]
						if delCron != nil {
							delCron.Stop()
							reply(rm, "本群自动报时已关闭")
						}
					})
				} else if strings.Index(rm.Msg.Text, "天气 ") == 0 {//获取天气情况
					weather := getWeather(strings.Replace(rm.Msg.Text, "天气 ", "", 1))
					if weather != "" {
						reply(rm, weather)
					}
				} else {
					//测试回复发送的消息
					switch rm.Type {
					case "PrivateMsg":
						sendPrivateMsg(rm.LogonQQ, rm.FromQQ.UIN, rm.Msg.Text)
						break
						//case "GroupMsg":
						//	sendGroupMsg(rm.LogonQQ, rm.FromGroup.GIN, rm.Msg.Text)
						//	break
					}
				}

			}()
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				logApi.Debug("write:", err)
				return
			}
		case <-interrupt:
			logApi.Debug("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logApi.Error("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

type ReceiveMessage struct {
	Type   string `json:"Type"`
	FromQQ struct {
		UIN      int    `json:"UIN"`
		NickName string `json:"NickName"`
	} `json:"FromQQ"`
	LogonQQ   int `json:"LogonQQ"`
	TimeStamp struct {
		Recv int `json:"Recv"`
		Send int `json:"Send"`
	} `json:"TimeStamp"`
	FromGroup struct {
		GIN int `json:"GIN"`
	} `json:"FromGroup"`
	Msg struct {
		Req         int    `json:"Req"`
		Seq         int64  `json:"Seq"`
		Type        int    `json:"Type"`
		SubType     int    `json:"SubType"`
		SubTempType int    `json:"SubTempType"`
		Text        string `json:"Text"`
		BubbleID    int    `json:"BubbleID"`
	} `json:"Msg"`
	Hb struct {
		Type int `json:"Type"`
	} `json:"Hb"`
	File struct {
		ID   string `json:"ID"`
		MD5  string `json:"MD5"`
		Name string `json:"Name"`
		Size int    `json:"Size"`
	} `json:"File"`
}
