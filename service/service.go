package service

import (
	"context"
	etcd "go.etcd.io/etcd/client/v3"
	"puhser/config"
	"strconv"
	"time"
)

func Init(c config.Config) *etcd.Client {
	var err error
	EClient, err := etcd.New(etcd.Config{
		Endpoints:   c.Etcd.EndPoints,
		DialTimeout: time.Duration(c.Etcd.DialTimeout) * time.Second,
	})
	if err != nil {
		panic(err.Error())
	}
	_, err = RegisterService(c, EClient)
	if err != nil {
		panic(err)
	}
	return EClient
}

func RegisterService(c config.Config, EClient *etcd.Client) (etcd.LeaseID, error) {
	grantResp, err := EClient.Grant(context.Background(), c.Etcd.TTL)
	if err != nil {
		return 0, err
	}
	name := c.Etcd.Name + "/" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	_, err = EClient.Put(context.Background(), name, c.Etcd.Addr, etcd.WithLease(grantResp.ID))
	if err != nil {
		return 0, err
	}
	ch, err := EClient.KeepAlive(context.Background(), grantResp.ID)
	if err != nil {
		return 0, err
	}
	go func() {
		for {
			select {
			case resp := <-ch:
				if resp == nil {
					panic("grant timeout")
				}
			}
		}
	}()
	return grantResp.ID, nil
}
