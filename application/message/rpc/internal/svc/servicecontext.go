package svc

import (
	"lifememo/application/message/rpc/internal/config"
	"lifememo/application/message/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	MessageModel model.MessageModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		MessageModel: model.NewMessageModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
	}
}
