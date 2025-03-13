package main

import (
	svc "puhser/internal/context"
	"puhser/mq"
	"puhser/route"
	"puhser/service"
)

func main() {
	ctx := svc.NewContext("./internal/config/config.yaml")

	go service.Init(ctx)
	mq.Init(ctx)
	route.Init(ctx)
}
