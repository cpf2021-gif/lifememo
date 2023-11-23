package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	DataSource string
	CacheRedis cache.CacheConf
	BizRedis   redis.RedisConf
}
