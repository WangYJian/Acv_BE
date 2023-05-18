package webSocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var Lock = sync.RWMutex{}

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// IdConnMap id和连接的映射
var IdConnMap = make(map[string]*websocket.Conn)
