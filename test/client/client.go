package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

const (
	url       = "ws://127.0.0.1:8077/ToWebSocket" // 固定 WebSocket 服务器地址
	connCount = 100                               // 并发连接数
)

var cnt = 0
var mu sync.Mutex

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(url, http.Header{"cookie": {"session_id=112233"}})
	if err != nil {
		panic(err.Error())
		return
	}
	mu.Lock()
	cnt++
	mu.Unlock()
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}()
	for {
		if err := conn.WriteMessage(websocket.TextMessage, []byte{}); err != nil {
			mu.Lock()
			cnt--
			mu.Unlock()
			return
		}
		time.Sleep(time.Second * 15)
	}

}
