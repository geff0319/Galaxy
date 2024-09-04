package bridge

import (
	"context"
	"encoding/json"
	"errors"
	"galaxy/bridge/ytdlp"
	"github.com/coder/websocket"
	"github.com/ge-fei-fan/gefflog"
	"net"
	"time"
)

type wsClent struct {
	id     string
	domain string
	wsConn *websocket.Conn
	close  chan struct{}
}

var WsC = wsClent{
	wsConn: nil,
}

func getWsUrl(domain, id string) string {
	return "ws://" + domain + "/subscribe/" + id
}
func resetWsClient() {
	WsC.domain = ""
	WsC.id = ""
	//WsC.wsConn = nil
	if WsC.close != nil {
		WsC.close <- struct{}{}
	}
}
func connect(domain, id string) (*websocket.Conn, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

	c, _, err := websocket.Dial(ctx, getWsUrl(domain, id), nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}
func (a *App) ConnectWs(domain, id string) FlagResult {
	if WsC.wsConn != nil {
		WsC.wsConn.CloseNow()
	}
	c, err := connect(domain, id)
	if err != nil {
		gefflog.Err("connect err: " + getWsUrl(domain, id))
		return FlagResult{Flag: false, Data: "ws连接失败"}
	}
	WsC.domain = domain
	WsC.id = id
	WsC.wsConn = c
	WsC.close = make(chan struct{}, 1)

	go nextMessage()
	go pingLoop()
	return FlagResult{Flag: true, Data: "ws连接成功"}
}
func (a *App) DisConnectWs() FlagResult {
	if WsC.wsConn != nil {
		err := WsC.wsConn.CloseNow()
		if err != nil {
			gefflog.Err("disconnect err")
			return FlagResult{Flag: false, Data: "ws关闭失败"}
		}
	}
	resetWsClient()
	return FlagResult{Flag: true, Data: "ws关闭成功"}
}
func (a *App) Ping(domain, id string) FlagResult {
	if WsC.wsConn == nil {
		return FlagResult{Flag: false, Data: "ws连接不存在"}
	}
	if WsC.domain != domain || WsC.id != id {
		WsC.wsConn.CloseNow()
		resetWsClient()
		return FlagResult{Flag: false, Data: "ws连接和当前配置不匹配"}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := WsC.wsConn.Ping(ctx)
	if err != nil {
		return FlagResult{Flag: false, Data: "ws连接已断开"}
	}
	return FlagResult{Flag: true, Data: "ws连接正常"}
}

type revMsg struct {
	Data string `json:"data"`
}

func nextMessage() {
	var rm revMsg
	for {
		_, b, err := WsC.wsConn.Read(context.Background())
		if err != nil {
			gefflog.Err(WsC.id + "获取WS内容报错：" + err.Error())
			if errors.Is(err, net.ErrClosed) {
				gefflog.Info("WS连接关闭退出")
				break
			}
			continue
		}
		err = json.Unmarshal(b, &rm)
		if err != nil {
			gefflog.Err("failed to unmarshal JSON: %w", err)
		}
		gefflog.Info("WS接收到消息" + rm.Data)
		url, ok := ytdlp.AppYoutubeCompile(rm.Data)
		if ok {
			p := &ytdlp.Process{
				Url:    url,
				Params: []string{},
				Output: ytdlp.DownloadOutput{},
			}
			ytdlp.YdpConfig.Mdb.Set(p)
			ytdlp.YdpConfig.Mq.Publish(p)
			continue
		}
		url, ok = ytdlp.AppXCompile(rm.Data)
		if ok {
			p := &ytdlp.Process{
				Url:    url,
				Params: []string{},
				Output: ytdlp.DownloadOutput{},
			}
			ytdlp.YdpConfig.Mdb.Set(p)
			ytdlp.YdpConfig.Mq.Publish(p)
			continue
		}
	}
}

func pingLoop() {
	ticker := time.NewTicker(time.Minute) // 每 2 秒发送一次 ping
	defer ticker.Stop()

	for {
		select {
		case <-WsC.close:
			// 上下文取消，退出 goroutine
			gefflog.Info("上下文取消，退出 ping goroutine")
			return
		case <-ticker.C:
			// 发送 ping 消息
			if err := WsC.wsConn.Ping(context.Background()); err != nil {
				gefflog.Err("发送 ping 消息失败:", err)
				return
			}
		}
	}
}
