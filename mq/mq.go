package mq

import (
	"github.com/nsqio/go-nsq"
)

type MessageHandler struct{}

type Message struct {
}

func (h *MessageHandler) HandleMessage(message *nsq.Message) error {
}
