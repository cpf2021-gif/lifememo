package svc

import (
	"lifememo/application/moment/rpc/internal/config"
	"lifememo/application/moment/rpc/internal/model"

	"github.com/golang/groupcache/singleflight"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	MomentModel       model.MomentModel
	BizRedis          *redis.Redis
	SingleFlightGroup singleflight.Group
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:      c,
		MomentModel: model.NewMomentModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		BizRedis:    rds,
	}
}
