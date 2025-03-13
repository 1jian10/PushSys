package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"puhser/client"
	"puhser/config"
)

var UpGrader websocket.Upgrader

func Init(c config.Config) {
	UpGrader = websocket.Upgrader{
		ReadBufferSize:  c.Websocket.ReadBufferSize,
		WriteBufferSize: c.Websocket.WriteBufferSize,
	}
	//gin.SetMode(gin.ReleaseMode)
	//gin.DefaultWriter = ioutil.Discard
	engine := gin.Default()
	engine.GET("/ToWebSocket", ToWebSocket)
	err := engine.Run(":" + c.Websocket.Port)
	if err != nil {
		panic(err)
	}
}

func ToWebSocket(ctx *gin.Context) {
	if !websocket.IsWebSocketUpgrade(ctx.Request) {
		return
	}
	session, err := ctx.Cookie("session_id")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conn, err := UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	client.NewClient(conn, session)
}
