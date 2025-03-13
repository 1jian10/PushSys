package svc

import (
	"InnerGetWay/internal/config"
	"github.com/redis/go-redis/v9"
	etcd "go.etcd.io/etcd/test/v3"
)

type ServiceContext struct {
	Config  config.Config
	RDB     *redis.Client
	EClient *etcd.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
