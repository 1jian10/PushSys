package service

import (
	svc "OutGetWay/context"
	"context"
	etcd "go.etcd.io/etcd/client/v3"
	"math/rand"
	"sync"
)

var OldAddr []string
var NewAddr []string
var Addr = make(map[string]string)
var mu sync.RWMutex

func Watch(ctx *svc.Context) {
	watcher := etcd.NewWatcher(ctx.EClient)
	WatchChan := watcher.Watch(context.Background(), ctx.Config.Etcd.WatchPrefix, etcd.WithPrefix())
	for resp := range WatchChan {
		for _, ev := range resp.Events {
			if ev.Type == etcd.EventTypePut {
				Addr[string(ev.Kv.Value)] = "1"
			} else {
				delete(Addr, string(ev.Kv.Value))
			}
		}
		for k := range Addr {
			NewAddr = append(NewAddr, k)
		}
		mu.Lock()
		OldAddr = NewAddr
		mu.Unlock()
	}
}

func InitService(ctx *svc.Context) {
	kv := etcd.NewKV(ctx.EClient)
	resp, err := kv.Get(context.Background(), ctx.Config.Etcd.WatchPrefix, etcd.WithPrefix())
	if err != nil {
		panic(err)
	}
	for _, v := range resp.Kvs {
		OldAddr = append(OldAddr, string(v.Value))
	}
	go Watch(ctx)
}

func SelectService() string {
	mu.RLock()
	defer mu.RUnlock()
	return OldAddr[rand.Intn(len(OldAddr))]
}
