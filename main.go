package main

import (
	"puhser/client"
	"puhser/config"
	rpc "puhser/proto"
	"puhser/route"
	"puhser/service"
)

func main() {
	c := config.ReadConfig()

	client.Init(c)
	service.Init(c)
	go rpc.Init(c)
	route.Init(c)

}
