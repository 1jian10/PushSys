package main

import (
	svc "OutGetWay/context"
	"OutGetWay/route"
	"fmt"
)

func main() {
	ctx := svc.NewContext("./config/config.yaml")
	fmt.Println(ctx.Config)
	route.Init(ctx)
}
