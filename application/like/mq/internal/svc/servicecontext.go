package svc

import (
	"lifememo/application/like/mq/internal/config"
	"lifememo/application/like/mq/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config      config.Config
	MomentModel model.MomentModel
	ReplyModel  model.ReplyModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		MomentModel: model.NewMomentModel(sqlx.NewMysql(c.MomentDatasource), c.CacheRedis),
		ReplyModel:  model.NewReplyModel(sqlx.NewMysql(c.ReplyDatasource), c.CacheRedis),
	}
}
