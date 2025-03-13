package mq

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nsqio/go-nsq"
	svc "puhser/internal/context"
	"puhser/route"
)

func Init(ctx *svc.Context) {
	config := nsq.NewConfig()
	c, err := nsq.NewConsumer(ctx.Config.NSQ.Topic, uuid.New().String(), config)
	if err != nil {
		panic(err.Error())
	}
	c.AddHandler(&MessageHandler{ctx: ctx})
	if err := c.ConnectToNSQD(ctx.Config.NSQ.Addr); err != nil {
		panic(err.Error())
	}
}

type MessageHandler struct {
	ctx *svc.Context
}
type Request struct {
	Msg      route.Message `json:"message"`
	BucketId int64         `json:"bucketId"`
}

func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	req := new(Request)
	if err := json.Unmarshal(m.Body, req); err != nil {
		m.Finish()
		return err
	}
	route.SendGlobalMessage(req.BucketId, &req.Msg)
	m.Finish()
	return nil
}
