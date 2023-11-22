package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	KqPusherConf struct {
		Brokers []string
		Topic   string
	}
	DataSource string
	CacheRedis cache.CacheConf
	BizRedis   redis.RedisConf
}
