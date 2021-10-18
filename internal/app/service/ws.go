package service

import (
	"github.com/LittleBenx86/Benlog/internal/global/consts"
	"github.com/LittleBenx86/Benlog/internal/global/variables"
	gin2 "github.com/LittleBenx86/Benlog/internal/utils/websocket/ginws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WebSocket struct {
	Client *gin2.Client
}

func (w *WebSocket) OnOpen(ctx *gin.Context) (*WebSocket, bool) {
	c, ok := (&gin2.Client{
		ClientContext: &gin2.ClientContext{
			ReadBufferSize:        variables.YmlAppConfig.GetInt("Websocket.RwBufferSize"),
			WriteBufferSize:       variables.YmlAppConfig.GetInt("Websocket.RwBufferSize"),
			PingPeriodSecond:      variables.YmlAppConfig.GetDuration("Websocket.PingPeriodSecond"),
			ReadDeadlineSecond:    variables.YmlAppConfig.GetDuration("Websocket.ReadDeadlineSecond"),
			WriteDeadlineSecond:   variables.YmlAppConfig.GetDuration("Websocket.WriteDeadlineSecond"),
			MaxMsgBytesSize:       variables.YmlAppConfig.GetInt64("Websocket.MaxMsgBytesSize"),
			MaxHeartbeatFailCount: variables.YmlAppConfig.GetInt("Websocket.MaxHeartbeatFailCount"),
			Logger:                variables.Logger,
		},
	}).OnOpen(ctx)
	if !ok {
		return nil, false
	}

	w.Client = c
	go w.Client.Heartbeat()
	return w, true
}

func (w *WebSocket) OnMessage(ctx *gin.Context) {
	go w.Client.Read(func(msgType int, receiveData []byte) {
		responseMsg := "server receive: " + string(receiveData)
		if err := w.Client.Send(msgType, responseMsg); err != nil {
			variables.Logger.Error(consts.ERRORS_WS_WRITE_MESSAGE_FAILED, zap.Error(err))
		}
	}, w.OnError, w.OnClose)
}

func (w *WebSocket) OnError(err error) {
	w.Client.State = gin2.INVALID // heartbeat will exit automatically
	variables.Logger.Error("client offline or suspend, browser refresh", zap.Error(err))
}

func (w *WebSocket) OnClose() {
	gin2.GetHubInstance().OfflineUnregister <- w.Client
}

func (w *WebSocket) GetOnlineClients() int {
	length := 0
	gin2.GetHubInstance().ClientCache.Range(func(k, v interface{}) bool {
		length++
		return true
	})
	return length
}

func (w *WebSocket) Broadcast(msg string) {
	gin2.GetHubInstance().ClientCache.Range(func(k, v interface{}) bool {
		c, ok := k.(*gin2.Client)
		if !ok {
			return false
		}
		_, ok = v.(bool)
		if !ok {
			return false
		}

		if err := c.Send(websocket.TextMessage, msg); err != nil {
			variables.Logger.Error(consts.ERRORS_WS_WRITE_MESSAGE_FAILED, zap.Error(err))
		}
		return false
	})
}
