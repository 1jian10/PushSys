package main

import (
	svc "OutGetWay/context"
	"OutGetWay/route"
)

// 外部网关，随机获取一个内部pusher节点地址
func main() {
	ctx := svc.NewContext("./config/config.yaml")
	route.Init(ctx)
}
