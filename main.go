package main

import (
	"bitbucket.org/Limard/logx"
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"
)

var addr = flag.String("addr", "192.168.58.134:10429", "http service address")

var (
	LoginQQ  int //当前登录QQ
	logApi   = logx.New(".", "xiaolizi-http")
	pongTime = 5 * time.Second
	pingTime = 5 * time.Second
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
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
				//msg := fmt.Sprintf("%s", message)
				//logApi.Debug(msg)
				var rm = new(ReceiveMessage)
				err = json.Unmarshal(message, &rm)
				if err != nil {
					logApi.Debug("转实体类出错")
				}
				//干其他事情
				logApi.Debugf("type: %v", rm.Type)
				//忽略自己的消息
				if rm.FromQQ.UIN == LoginQQ {
					return
				}

				//获取天气情况
				if strings.Index(rm.Msg.Text, "天气 ") == 0 {
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
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
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
