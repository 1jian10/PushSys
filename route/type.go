package route

import (
	"github.com/gorilla/websocket"
	svc "puhser/internal/context"
	"sync"
)

type Client struct {
	Session  string
	userId   string
	bucketId int64
	conn     *websocket.Conn
}

type Message struct {
	UUId       string `json:"uuid"`
	Session    string `json:"session"`
	PayLoad    string `json:"payLoad"`
	EncodeType string `json:"encodeType"`
}

var Ctx *svc.Context
var UpGrader websocket.Upgrader
var Bucket = make([]sync.Map, 100)
