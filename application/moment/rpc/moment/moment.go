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
	MomentDeleteRequest  = pb.MomentDeleteRequest
	MomentDeleteResponse = pb.MomentDeleteResponse
	MomentDetailRequest  = pb.MomentDetailRequest
	MomentDetailResponse = pb.MomentDetailResponse
	MomentItem           = pb.MomentItem
	MomentsRequest       = pb.MomentsRequest
	MomentsResponse      = pb.MomentsResponse
	PublishRequest       = pb.PublishRequest
	PublishResponse      = pb.PublishResponse

	Moment interface {
		Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error)
		Moments(ctx context.Context, in *MomentsRequest, opts ...grpc.CallOption) (*MomentsResponse, error)
		MomentDelete(ctx context.Context, in *MomentDeleteRequest, opts ...grpc.CallOption) (*MomentDeleteResponse, error)
		MomentDetail(ctx context.Context, in *MomentDetailRequest, opts ...grpc.CallOption) (*MomentDetailResponse, error)
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

func (m *defaultMoment) Moments(ctx context.Context, in *MomentsRequest, opts ...grpc.CallOption) (*MomentsResponse, error) {
	client := pb.NewMomentClient(m.cli.Conn())
	return client.Moments(ctx, in, opts...)
}

func (m *defaultMoment) MomentDelete(ctx context.Context, in *MomentDeleteRequest, opts ...grpc.CallOption) (*MomentDeleteResponse, error) {
	client := pb.NewMomentClient(m.cli.Conn())
	return client.MomentDelete(ctx, in, opts...)
}

func (m *defaultMoment) MomentDetail(ctx context.Context, in *MomentDetailRequest, opts ...grpc.CallOption) (*MomentDetailResponse, error) {
	client := pb.NewMomentClient(m.cli.Conn())
	return client.MomentDetail(ctx, in, opts...)
}
