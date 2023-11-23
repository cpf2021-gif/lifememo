package svc

import (
	"lifememo/application/like/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config   config.Config
	BizRedis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	res, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:   c,
		BizRedis: res,
	}
}
