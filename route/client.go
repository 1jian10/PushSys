package route

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"strconv"
	"time"
)

func (c *Client) Send(message *Message) error {
	if err := c.conn.SetWriteDeadline(time.Now().Add(time.Second * 5)); err != nil {
		return err
	}
	msg, err := json.Marshal(*message)
	if err != nil {
		return err
	}
	if err = c.conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
		return err
	}
	return nil
}

func GetUserId(session string) (int64, error) {
	return rand.Int63(), nil
}

func NewClient(conn *websocket.Conn, session string) {
	rdb := Ctx.RDB
	id, err := GetUserId(session)
	if err != nil {
		fmt.Println(err.Error())
	}
	c := &Client{
		Session:  session,
		userId:   strconv.FormatInt(id, 10),
		bucketId: id % 100,
		conn:     conn,
	}
	rdb.Set(context.Background(), "pusher:"+c.userId, Ctx.Config.Etcd.Addr, 2*time.Minute)
	Bucket[c.bucketId].Store(c.userId, c)
	go c.HeartCheck()
}

func (c *Client) HeartCheck() {
	rdb := Ctx.RDB
	for {
		_ = c.conn.SetReadDeadline(time.Now().Add(time.Minute))
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			c.CloseConn()
			return
		}
		rdb.Expire(context.Background(), "pusher:"+c.userId, 2*time.Minute)
	}
}

func (c *Client) CloseConn() {
	Ctx.RDB.Del(context.Background(), "pusher:"+c.userId)
	Bucket[c.bucketId].Delete(c.userId)
	_ = c.conn.Close()
	return
}
