package main

import (
	"context"
	"github.com/gin-gonic/gin"
	etcd "go.etcd.io/etcd/test/v3"
	"time"
)

func main() {
	EClient, err := etcd.New(etcd.Config{
		Endpoints:   []string{"127.0.0.1:4379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	engine := gin.Default()
	engine.GET("/GetHosts")
	engine.Run()
}

func GetHost(c *gin.Context) {

}

func WatchServices(EClient *etcd.Client) {
	watcher := etcd.NewWatcher(EClient)
	watchChan := watcher.Watch(context.Background(), "pusher", etcd.WithPrefix())
	for resp := range watchChan {
		for _, ev := range resp.Events {
			if ev.Type == etcd.EventTypePut {

			} else {

			}
		}
	}
}
