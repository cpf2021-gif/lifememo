package svc

import (
	"lifememo/application/like/rpc/internal/config"
	"lifememo/application/like/rpc/internal/model"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                   config.Config
	BizRedis                 *redis.Redis
	LikeRecordModel          model.LikeRecordModel
	LikeCountModel           model.LikeCountModel
	MomentLikeKqPusherClient *kq.Pusher
	ReplyLikeKqPusherClient  *kq.Pusher
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
		Config:                   c,
		BizRedis:                 res,
		LikeRecordModel:          model.NewLikeRecordModel(conn, c.CacheRedis),
		LikeCountModel:           model.NewLikeCountModel(conn, c.CacheRedis),
		MomentLikeKqPusherClient: kq.NewPusher(c.MomentLikeKqPusherConf.Brokers, c.MomentLikeKqPusherConf.Topic),
		ReplyLikeKqPusherClient:  kq.NewPusher(c.ReplyLikeKqPusherConf.Brokers, c.ReplyLikeKqPusherConf.Topic),
	}
}
