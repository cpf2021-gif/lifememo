package svc

import (
	"lifememo/application/follow/api/internal/config"
	"lifememo/application/follow/rpc/follow"
	"lifememo/application/user/rpc/user"
	"lifememo/pkg/interceptors"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	FollowRPC follow.Follow
	UserRPC   user.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	followRPC := zrpc.MustNewClient(c.FollowRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	userRPC := zrpc.MustNewClient(c.UserRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:    c,
		FollowRPC: follow.NewFollow(followRPC),
		UserRPC:   user.NewUser(userRPC),
	}
}
