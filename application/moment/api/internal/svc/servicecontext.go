package svc

import (
	"lifememo/application/moment/api/internal/config"
	"lifememo/application/moment/rpc/moment"
	"lifememo/pkg/interceptors"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	MomentRPC moment.Moment
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	momentRPC := zrpc.MustNewClient(c.MomentRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:    c,
		MomentRPC: moment.NewMoment(momentRPC),
	}
}
