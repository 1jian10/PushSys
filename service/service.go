package service

import (
	"context"
	"fmt"
	etcd "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"net"
	svc "puhser/internal/context"
	"puhser/internal/server"
	"puhser/proto/push"
	"strconv"
	"time"
)

func Init(ctx *svc.Context) {
	var err error

	_, err = RegisterService(ctx)
	if err != nil {
		panic(err)
	}
	lis, err := net.Listen("tcp", ctx.Config.Etcd.Addr)
	if err != nil {
		panic(err.Error())
	}
	grpcServer := grpc.NewServer()
	push.RegisterPushMessageServiceServer(grpcServer, server.NewPushMessageServiceServer(ctx))
	fmt.Println("rpc begin listening:" + ctx.Config.Etcd.Addr)
	if err = grpcServer.Serve(lis); err != nil {
		panic(err.Error())
	}
}

func RegisterService(ctx *svc.Context) (etcd.LeaseID, error) {
	EClient := ctx.EClient
	c := ctx.Config
	grantResp, err := EClient.Grant(context.Background(), c.Etcd.TTL)
	if err != nil {
		return 0, err
	}
	name := c.Etcd.Name + "/" + strconv.FormatInt(time.Now().UnixNano(), 10)
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
