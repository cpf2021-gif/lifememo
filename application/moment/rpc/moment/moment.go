// Code generated by goctl. DO NOT EDIT.
// Source: moment.proto

package moment

import (
	"context"

	"lifememo/application/moment/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	PublishRequest  = pb.PublishRequest
	PublishResponse = pb.PublishResponse

	Moment interface {
		Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error)
	}

	defaultMoment struct {
		cli zrpc.Client
	}
)

func NewMoment(cli zrpc.Client) Moment {
	return &defaultMoment{
		cli: cli,
	}
}

func (m *defaultMoment) Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error) {
	client := pb.NewMomentClient(m.cli.Conn())
	return client.Publish(ctx, in, opts...)
}
