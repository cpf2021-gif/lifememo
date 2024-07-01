package svc

import (
	"lifememo/application/like/api/internal/config"
	"lifememo/application/like/rpc/like"
	"lifememo/application/moment/rpc/moment"
	"lifememo/application/reply/rpc/reply"
	"lifememo/pkg/interceptors"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	LikeRPC   like.Like
	MomentRPC moment.Moment
	ReplyRPC  reply.Reply
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	likeRPC := zrpc.MustNewClient(c.LikeRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	momentRPC := zrpc.MustNewClient(c.MomentRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	replyRPC := zrpc.MustNewClient(c.ReplyRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:    c,
		LikeRPC:   like.NewLike(likeRPC),
		MomentRPC: moment.NewMoment(momentRPC),
		ReplyRPC:  reply.NewReply(replyRPC),
	}
}
