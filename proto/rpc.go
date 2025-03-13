package rpc

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"puhser/config"
	"puhser/internal/server"
	"puhser/internal/svc"
	"puhser/proto/push"
)

func Init(c config.Config) {
	lis, err := net.Listen("tcp", c.Etcd.Addr)
	if err != nil {
		panic(err.Error())
	}
	grpcServer := grpc.NewServer()
	ctx := svc.NewServiceContext(c)
	push.RegisterPushMessageServiceServer(grpcServer, server.NewPushMessageServiceServer(ctx))
	fmt.Println("rpc begin listening:" + c.Etcd.Addr)
	if err = grpcServer.Serve(lis); err != nil {
		panic(err.Error())
	}

}
