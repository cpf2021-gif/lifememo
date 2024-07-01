package svc

import (
	"lifememo/application/message/api/internal/config"
	"lifememo/application/message/rpc/messageservice"
	"lifememo/application/user/rpc/user"
	"lifememo/pkg/interceptors"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	MessageRPC messageservice.MessageService
	UserRPC    user.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	messageRPC := zrpc.MustNewClient(c.MessageRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	userRPC := zrpc.MustNewClient(c.UserRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:     c,
		MessageRPC: messageservice.NewMessageService(messageRPC),
		UserRPC:    user.NewUser(userRPC),
	}
}
