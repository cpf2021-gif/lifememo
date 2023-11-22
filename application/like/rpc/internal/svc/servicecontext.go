package svc

import (
	"lifememo/application/like/rpc/internal/config"
	"lifememo/application/like/rpc/internal/model"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	KqPusherClient  *kq.Pusher
	LikeRecordModel model.LikeRecordModel
	LikeCountModel  model.LikeCountModel
	BizRedis        *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)

	res, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:          c,
		KqPusherClient:  kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
		LikeRecordModel: model.NewLikeRecordModel(conn, c.CacheRedis),
		LikeCountModel:  model.NewLikeCountModel(conn, c.CacheRedis),
		BizRedis:        res,
	}
}
