package svc

import (
	"lifememo/application/reply/rpc/internal/config"
	"lifememo/application/reply/rpc/internal/model"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                      config.Config
	ReplyModel                  model.ReplyModel
	ReplyCountModel             model.ReplyCountModel
	BizRedis                    *redis.Redis
	MomentCommentKqPusherClient *kq.Pusher
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
		Config:                      c,
		ReplyModel:                  model.NewReplyModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		ReplyCountModel:             model.NewReplyCountModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		BizRedis:                    rds,
		MomentCommentKqPusherClient: kq.NewPusher(c.MomentCommentKqPusherConf.Brokers, c.MomentCommentKqPusherConf.Topic),
	}
}
