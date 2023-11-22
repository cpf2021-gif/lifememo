package svc

import (
	"lifememo/application/like/mq/internal/config"
	"lifememo/application/like/mq/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	LikeRecordModel model.LikeRecordModel
	LikeCountModel  model.LikeCountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)

	return &ServiceContext{
		Config:          c,
		LikeRecordModel: model.NewLikeRecordModel(conn, c.CacheRedis),
		LikeCountModel:  model.NewLikeCountModel(conn, c.CacheRedis),
	}
}
