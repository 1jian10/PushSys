package main

import (
	svc "puhser/internal/context"
	"puhser/route"
	"puhser/service"
)

func main() {
	ctx := svc.NewContext("./internal/config/config.yaml")

	go service.Init(ctx)
	route.Init(ctx)

}
