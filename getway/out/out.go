package main

import (
	svc "OutGetWay/context"
	"OutGetWay/route"
)

func main() {
	ctx := svc.NewContext("./config/config.yaml")
	route.Init(ctx)
}
