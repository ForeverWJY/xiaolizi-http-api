package main

import (
	"bitbucket.org/Limard/logx"
	"encoding/json"
	"fmt"
	"github.com/coocood/freecache"
	"github.com/gorilla/websocket"
	"github.com/robfig/cron"
	"net/http"
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
	logApi  = logx.New("./logs", "xiaolizi-http")
	//pongTime = 5 * time.Second
	//pingTime = 5 * time.Second
	cacheSize = 10 * 1024 * 1024 // In bytes, where 1024 * 1024 represents a single Megabyte, and 10 * 1024*1024 represents 10 Megabytes.
	cache     = freecache.NewCache(cacheSize)
	addr      string
	httpAddr  string
	managerQQ = make(map[int]int)
	//群号 : cron
	autoTimeMap = make(map[int]*cron.Cron)
	jsonConfig  Config
)

// Find获取一个切片并在其中查找元素。如果找到它，它将返回它的密钥，否则它将返回-1和一个错误的bool。
func FindSlice(slice []int, val int) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func loadManager(qqList []int) {
	for _, i := range qqList {
		fmt.Println(i)
		managerQQ[i] = i
	}
}

func isMangerQQ(rm *EventResponse, f func()) {
	i := managerQQ[rm.UserID]
	if i > 0 {
		if f != nil {
			f()
		}
	} else {
		reply(rm, "仅管理员可操作")
	}
}

func loadAutoTimeCron() {
	for _, v := range jsonConfig.ReportTime.Group {
		resp := new(EventResponse)
		resp.MessageType = "group"
		resp.GroupID = v
		newCron := cron.New()
		err2 := newCron.AddFunc("0 0 * * * *", func() {
			reply(resp, time.Now().Format("2006-01-02 15:04:05"))
		})
		if err2 != nil {
			logx.Error(err2.Error())
		}
		newCron.Start()
		autoTimeMap[v] = newCron
	}
	for _, v := range jsonConfig.ReportTime.Qq {
		resp := new(EventResponse)
		resp.MessageType = "private"
		resp.Sender.UserID = v
		newCron := cron.New()
		err2 := newCron.AddFunc("0 0 * * * *", func() {
			reply(resp, time.Now().Format("2006-01-02 15:04:05"))
		})
		if err2 != nil {
			logx.Error(err2.Error())
		}
		newCron.Start()
		autoTimeMap[v] = newCron
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "go-cqhttp-api golang")
}

type EventResponse struct {
	Interval      int    `json:"interval"`
	MetaEventType string `json:"meta_event_type"`
	PostType      string `json:"post_type"` //事件类型
	SelfID        int    `json:"self_id"`   //收到事件的机器人 QQ 号
	Status        struct {
		AppEnabled     bool        `json:"app_enabled"`
		AppGood        bool        `json:"app_good"`
		AppInitialized bool        `json:"app_initialized"`
		Good           bool        `json:"good"`
		Online         bool        `json:"online"`
		PluginsGood    interface{} `json:"plugins_good"`
		Stat           struct {
			PacketReceived  int `json:"packet_received"`
			PacketSent      int `json:"packet_sent"`
			PacketLost      int `json:"packet_lost"`
			MessageReceived int `json:"message_received"`
			MessageSent     int `json:"message_sent"`
			DisconnectTimes int `json:"disconnect_times"`
			LostTimes       int `json:"lost_times"`
		} `json:"stat"`
	} `json:"status"`
	Sender struct {
		Age      int    `json:"age"`
		Nickname string `json:"nickname"`
		Sex      string `json:"sex"`
		UserID   int    `json:"user_id"`
	} `json:"sender"`
	SubType     string `json:"sub_type"`
	Time        int    `json:"time"`     //	事件发生的时间戳
	UserID      int    `json:"user_id"`  //发送者QQ
	GroupID     int    `json:"group_id"` //群号
	Font        int    `json:"font"`
	Message     string `json:"message"`
	MessageID   int    `json:"message_id"`
	MessageType string `json:"message_type"`
	RawMessage  string `json:"raw_message"`
}

func main() {
	//flag.Parse()
	//log.SetFlags(0)

	jsonConfig = loadJsonConfig()
	addr = fmt.Sprintf("%v:%v", jsonConfig.IP, jsonConfig.WsPort)
	httpAddr = fmt.Sprintf("%v:%v", jsonConfig.IP, jsonConfig.Port)

	if jsonConfig.ManagerQQ != nil {
		loadManager(jsonConfig.ManagerQQ)
	}
	loadAutoTimeCron()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: addr, Path: "/event"}
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
	//LoginQQ = getLoginQQ()

	testFunc()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logApi.Debug("read error :", err.Error())
				return
			}
			go func() {
				logApi.Debugf("recv: %s", message)
				var rm = new(EventResponse)
				err = json.Unmarshal(message, &rm)
				if err != nil {
					logApi.Debug("转实体类出错", err.Error())
					return
				}
				//得到当前QQ
				LoginQQ = rm.SelfID
				//干其他事情
				logApi.Debugf("post type: %v", rm.PostType)
				go onReceiveMessage(rm)
			}()
		}
	}()

	go func() {
		//http server
		http.HandleFunc("/", indexHandler)
		http.ListenAndServe(fmt.Sprintf(":%v", jsonConfig.HttpPort), nil)
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

func onReceiveMessage(rm *EventResponse) {
	switch rm.PostType {
	case "message": //消息事件
		break
	case "notice": //通知事件
		/*
			群文件上传
			群管理员变动
			群成员减少
			群成员增加
			群禁言
			好友添加
			群消息撤回
			好友消息撤回
			群内戳一戳
			群红包运气王
			群成员荣誉变更
		*/
		return
		break
	case "request": //请求事件
		return
		break
	case "meta_event": //元事件
		break
	}
	groupQQ := rm.GroupID
	if rm.Message == "help" {
		jsonConfig := loadJsonConfig()
		reply(rm, jsonConfig.Help)
	} else if rm.Message == "开启报时" {
		isMangerQQ(rm, func() {
			newCron := cron.New()
			err2 := newCron.AddFunc("0 0 * * * *", func() {
				reply(rm, time.Now().Format("2006-01-02 15:04:05"))
			})
			if err2 != nil {
				logx.Error(err2.Error())
			}
			newCron.Start()
			autoTimeMap[groupQQ] = newCron
			reply(rm, "自动报时已开启")

			change := false
			if groupQQ != 0 {
				_, found := FindSlice(jsonConfig.ReportTime.Group, rm.UserID)
				if !found {
					jsonConfig.ReportTime.Group = append(jsonConfig.ReportTime.Group, groupQQ)
					change = true
				}
			} else {
				_, found := FindSlice(jsonConfig.ReportTime.Qq, rm.UserID)
				if !found {
					jsonConfig.ReportTime.Qq = append(jsonConfig.ReportTime.Qq, rm.UserID)
					change = true
				}
			}
			if change {
				saveJsonConfig(jsonConfig)
			}
		})
	} else if rm.Message == "关闭报时" {
		isMangerQQ(rm, func() {
			delCron := autoTimeMap[groupQQ]
			if delCron != nil {
				delCron.Stop()
				reply(rm, "自动报时已关闭")

				//删除配置
				change := false
				if groupQQ != 0 {
					index, found := FindSlice(jsonConfig.ReportTime.Group, groupQQ)
					if found {
						// 将删除点前后的元素连接起来
						jsonConfig.ReportTime.Group = append(jsonConfig.ReportTime.Group[:index], jsonConfig.ReportTime.Group[index+1:]...)
						change = true
					}
				} else {
					index, found := FindSlice(jsonConfig.ReportTime.Qq, rm.UserID)
					if found {
						// 将删除点前后的元素连接起来
						jsonConfig.ReportTime.Qq = append(jsonConfig.ReportTime.Qq[:index], jsonConfig.ReportTime.Qq[index+1:]...)
						change = true
					}
				}
				if change {
					saveJsonConfig(jsonConfig)
				}
			}
		})
	} else if strings.Index(rm.Message, "天气 ") == 0 { //获取天气情况
		weather := getWeather(strings.Replace(rm.Message, "天气 ", "", 1))
		if weather != "" {
			reply(rm, weather)
		}
	} else if rm.Message == "oschina" {
		arr := getOSCHINA()
		if arr != nil {
			index := 0
			var rs []string
			var trs []string
			for _, v := range arr {
				index++
				trs = append(trs, string(rune(index))+"、"+v)
				//每10条发送一次消息
				if index%10 == 0 || index == len(arr) {
					rs = append(rs, strings.Join(trs, "\n"))
					trs = append([]string{})
				}
			}
			for _, i := range rs {
				reply(rm, i)
			}
			//reply(rm, strings.Join(arr, "\n"))
		}
	} else {
		//测试回复发送的消息
		switch rm.MessageType {
		case "private":
			sendPrivateMsg(rm.UserID, rm.Message)
			break
			//case "GroupMsg":
			//	sendGroupMsg(rm.LogonQQ, rm.FromGroup.GIN, rm.Msg.Text)
			//	break
		}
	}
}
