package svc

import (
	"lifememo/application/reply/mq/internal/config"
	"lifememo/application/reply/mq/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config      config.Config
	MomentModel model.MomentModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		MomentModel: model.NewMomentModel(sqlx.NewMysql(c.Datasource), c.CacheRedis),
	}
}
