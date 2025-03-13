package svc

import (
	"context"
	"github.com/nsqio/go-nsq"
	"github.com/redis/go-redis/v9"
	etcd "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"puhser/getway/in/internal/config"
	"puhser/proto/push"
	"sync"
	"time"
)

type ServiceContext struct {
	Config   config.Config
	RDB      *redis.Client
	Services sync.Map
	EClient  *etcd.Client
	Producer *nsq.Producer
}

func NewServiceContext(c config.Config) *ServiceContext {
	EClient, err := etcd.New(etcd.Config{
		Endpoints:   c.Etcd.Hosts,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err.Error())
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: c.MRedis.Host,
		DB:   c.MRedis.DB,
	})
	if err = rdb.Ping(context.Background()).Err(); err != nil {
		panic(err.Error())
	}
	producer, err := nsq.NewProducer(c.NSQ.Addr, nsq.NewConfig())
	if err != nil {
		panic(err.Error())
	}

	svc := &ServiceContext{
		Config:   c,
		EClient:  EClient,
		Producer: producer,
	}
	InitService(svc)
	go Watch(svc)

	return svc
}

func Watch(svc *ServiceContext) {
	watcher := etcd.NewWatcher(svc.EClient)
	WatchChan := watcher.Watch(context.Background(), svc.Config.WatchPrefix, etcd.WithPrefix())
	for resp := range WatchChan {
		for _, ev := range resp.Events {
			if ev.Type == etcd.EventTypePut {
				ConnService(svc, string(ev.Kv.Value))
			} else {
				svc.Services.Delete(string(ev.Kv.Value))
			}
		}
	}
}

func InitService(svc *ServiceContext) {
	kv := etcd.NewKV(svc.EClient)
	resp, err := kv.Get(context.Background(), svc.Config.WatchPrefix, etcd.WithPrefix())
	if err != nil {
		panic(err)
	}
	for _, v := range resp.Kvs {
		ConnService(svc, string(v.Value))
	}
}

func ConnService(svc *ServiceContext, addr string) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	c := push.NewPushMessageServiceClient(conn)
	svc.Services.Store(addr, c)
	return
}
