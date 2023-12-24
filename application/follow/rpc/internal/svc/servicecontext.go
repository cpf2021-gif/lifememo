package svc

import (
	"lifememo/application/follow/rpc/internal/config"
	"lifememo/application/follow/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config           config.Config
	FollowModel      model.FollowModel
	FollowCountModel model.FollowCountModel
	BizRedis         *redis.Redis
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
		Config:           c,
		FollowModel:      model.NewFollowModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		FollowCountModel: model.NewFollowCountModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		BizRedis:         rds,
	}
}
