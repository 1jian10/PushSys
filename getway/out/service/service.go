package service

import (
	svc "OutGetWay/context"
	"context"
	etcd "go.etcd.io/etcd/client/v3"
	"math/rand"
	"sync"
)

// OldAddr 存储旧pusher节点地址的列表，实时性不强
var OldAddr []string

// NewAddr 用于更新时的中间存储
var NewAddr []string

// Addr 存储节点的最新地址
var Addr = make(map[string]string)

// mu 读写锁，更新时+写锁，否则+读锁
var mu sync.RWMutex

// Watch 观测etcd中键值对的变化，对本地地址进行更新
func Watch(ctx *svc.Context) {
	watcher := etcd.NewWatcher(ctx.EClient)
	//根据前缀对key进行watch，返回值为一个channel，当有变化时，将会发送消息到该chan中
	WatchChan := watcher.Watch(context.Background(), ctx.Config.Etcd.WatchPrefix, etcd.WithPrefix())
	if ctx.Config.Model == 1 {
		randomWatch(WatchChan)
	} else {

	}
}

func randomWatch(WatchChan etcd.WatchChan) {
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
		//更新地址，因slice为引用，速度较快
		OldAddr = NewAddr
		mu.Unlock()
	}
}

func consistentHashWatch(WatchChan etcd.WatchChan) {
	for resp := range WatchChan {
		for _, ev := range resp.Events {
			if ev.Type == etcd.EventTypePut {

			} else {

			}
		}
		for k := range Addr {
			NewAddr = append(NewAddr, k)
		}
		mu.Lock()
		//更新地址，因slice为引用，速度较快
		OldAddr = NewAddr
		mu.Unlock()
	}

}

// InitService 初始化pusher节点列表
func InitService(ctx *svc.Context) {
	kv := etcd.NewKV(ctx.EClient)
	resp, err := kv.Get(context.Background(), ctx.Config.Etcd.WatchPrefix, etcd.WithPrefix())
	if err != nil {
		panic(err)
	}

	if ctx.Config.Model == 1 {
		for _, v := range resp.Kvs {
			OldAddr = append(OldAddr, string(v.Value))
		}
	}

	if ctx.Config.Model == 2 {

	}
	go Watch(ctx)
}

// SelectService 随机路由，选取一个随机的路由地址
func SelectService() string {
	mu.RLock()
	defer mu.RUnlock()
	return OldAddr[rand.Intn(len(OldAddr))]
}
