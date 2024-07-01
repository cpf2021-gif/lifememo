package svc

import (
	"lifememo/application/moment/rpc/moment"
	"lifememo/application/reply/api/internal/config"
	"lifememo/application/reply/rpc/reply"
	"lifememo/application/user/rpc/user"
	"lifememo/pkg/interceptors"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	ReplyRPC  reply.Reply
	MomentRPC moment.Moment
	UserRPC   user.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	replyRPC := zrpc.MustNewClient(c.ReplyRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	momentRPC := zrpc.MustNewClient(c.MomentRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	userRPC := zrpc.MustNewClient(c.UserRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:    c,
		ReplyRPC:  reply.NewReply(replyRPC),
		MomentRPC: moment.NewMoment(momentRPC),
		UserRPC:   user.NewUser(userRPC),
	}
}
