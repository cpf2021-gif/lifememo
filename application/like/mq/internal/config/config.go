package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/cache"
)

type Config struct {
	service.ServiceConf

	MomentLikeKqConsumerConf kq.KqConf
	ReplyLikeKqConsumerConf  kq.KqConf
	MomentDatasource         string
	ReplyDatasource          string
	CacheRedis               cache.CacheConf
}
