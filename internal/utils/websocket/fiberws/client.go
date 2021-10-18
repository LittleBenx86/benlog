package fiberws

import (
	SysErrors "errors"
	"sync"
	"time"

	"github.com/LittleBenx86/Benlog/internal/global/consts"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

/*
 Migrating from Gin to fiber web framework, the original websocket implemented by gorilla which not supported by fiber.
*/

type ClientContext struct {
	ReadBufferSize        int
	WriteBufferSize       int
	PingPeriodSecond      time.Duration
	ReadDeadlineSecond    time.Duration
	WriteDeadlineSecond   time.Duration
	MaxMsgBytesSize       int64
	MaxHeartbeatFailCount int
	Logger                *zap.Logger `json:"-" mapstructure:",remain"`
}

func (c *ClientContext) GetLogger() *zap.Logger {
	return c.Logger
}

type Client struct {
	Conn               *websocket.Conn
	SendCh             chan []byte
	PingPeriod         time.Duration
	ReadDeadline       time.Duration
	WriteDeadline      time.Duration
	HeartbeatFailCount int
	State              uint8 //0: error, 1: ok
	sync.RWMutex
	*ClientContext
}

const (
	STATE_OK uint8 = 1
	INVALID  uint8 = 0
)

func (c *Client) Send(msgType int, msg string) error {
	c.Lock()
	defer func() {
		c.Unlock()
	}()

	if err := c.Conn.SetWriteDeadline(time.Now().Add(c.WriteDeadline)); err != nil {
		c.GetLogger().Error(consts.ERRORS_WS_SET_WRITE_DEADLINE_FAILED, zap.Error(err))
		return err
	}

	if err := c.Conn.WriteMessage(msgType, []byte(msg)); err != nil {
		return err
	}
	return nil
}

func (c *Client) OnOpen(ctx *fiber.Ctx) (*Client, bool) {
	// upgrade protocol from http to websocket
	defer func() {
		err := recover()
		if err != nil {
			if val, ok := err.(error); ok {
				c.GetLogger().Error(consts.ERRORS_WS_OPEN_FAILED, zap.Error(val))
			}
		}
	}()

	upgrader := websocket.FastHTTPUpgrader{
		ReadBufferSize:  c.ReadBufferSize,
		WriteBufferSize: c.WriteBufferSize,
	}

	// init a long term active websocket client
	err := upgrader.Upgrade(ctx.Context(), func(conn *websocket.Conn) {
		c.Conn = conn
	})
	if err != nil {
		c.GetLogger().Error(consts.ERRORS_WS_PROTOCOL_UPGRADE_FAILED, zap.Error(err))
		return nil, false
	}

	c.SendCh = make(chan []byte, c.WriteBufferSize)
	c.PingPeriod = time.Second * c.PingPeriodSecond
	c.ReadDeadline = time.Second * c.ReadDeadlineSecond
	c.WriteDeadline = time.Second * c.WriteDeadlineSecond

	if err := c.Send(websocket.TextMessage, consts.WS_HANDSHAKE_SUCCEEDED); err != nil {
		c.GetLogger().Error(consts.ERRORS_WS_WRITE_MESSAGE_FAILED, zap.Error(err))
	}

	c.Conn.SetReadLimit(c.MaxMsgBytesSize)
	GetHubInstance().OnlineRegister <- c // register current client to hub
	c.State = STATE_OK
	return c, true
}

func (c *Client) Read(onMsgFn func(msgType int, data []byte), onErrFn func(err error), onCloseFn func()) {
	defer func() {
		err := recover()
		if err != nil {
			if e, ok := err.(error); ok {
				c.GetLogger().Error(consts.ERRORS_WS_READ_MESSAGE_FAILED, zap.Error(e))
			}
		}
		onCloseFn()
	}()

	for {
		if c.State == STATE_OK {
			message, receiveBytes, e := c.Conn.ReadMessage()
			if e != nil {
				onErrFn(e)
				break
			}
			onMsgFn(message, receiveBytes)
		} else if c.State == INVALID {
			onErrFn(SysErrors.New(consts.ERRORS_WS_STATE_INVALID))
			break
		}
	}
}

func (c *Client) Heartbeat() {
	ticker := time.NewTicker(c.PingPeriod)
	defer func() {
		err := recover()
		if err != nil {
			if e, ok := err.(error); ok {
				c.GetLogger().Error(consts.ERRORS_WS_HEARTBEAT_SERVER_ERROR, zap.Error(e))
			}
		}
		ticker.Stop()
	}()

	// browser receive ping will auto response pong
	if c.ReadDeadline == 0 {
		_ = c.Conn.SetReadDeadline(time.Time{})
	} else {
		_ = c.Conn.SetReadDeadline(time.Now().Add(c.ReadDeadline))
	}

	c.Conn.SetPongHandler(func(pong string) error {
		if c.ReadDeadline > time.Nanosecond {
			_ = c.Conn.SetReadDeadline(time.Now().Add(c.ReadDeadline))
		} else {
			_ = c.Conn.SetReadDeadline(time.Time{})
		}
		return nil
	})

	for {
		select {
		case <-ticker.C:
			if c.State == STATE_OK {
				if err := c.Send(websocket.PingMessage, consts.WS_SERVER_PING_MESSAGE); err != nil {
					c.HeartbeatFailCount++
					if c.HeartbeatFailCount > c.MaxHeartbeatFailCount {
						c.State = INVALID
						c.GetLogger().Error(consts.ERRORS_WS_HEARTBEAT_FAILURE_EXCEED_MAX, zap.Error(err))
						return
					}
				}

				if c.HeartbeatFailCount > 0 {
					c.HeartbeatFailCount--
				}
			} else if c.State == INVALID {
				return
			}
		}
	}
}
