package ginws

import "sync"

type Hub struct {
	OnlineRegister    chan *Client
	OfflineUnregister chan *Client
	ClientCache       sync.Map // key: *Client; value: bool
}

var (
	hubInstance *Hub
	initLock    sync.Mutex
)

func GetHubInstance() *Hub {
	if hubInstance == nil {
		initLock.Lock()
		defer initLock.Unlock()
		hubInstance = &Hub{
			OnlineRegister:    make(chan *Client),
			OfflineUnregister: make(chan *Client),
		}
	}
	return hubInstance
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.OnlineRegister:
			h.ClientCache.Store(c, true)
		case c := <-h.OfflineUnregister:
			if _, ok := h.ClientCache.Load(c); ok {
				_ = c.Conn.Close()
				h.ClientCache.Delete(c)
			}
		}
	}
}
