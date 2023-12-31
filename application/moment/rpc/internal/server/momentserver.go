// Code generated by goctl. DO NOT EDIT.
// Source: moment.proto

package server

import (
	"context"

	"lifememo/application/moment/rpc/internal/logic"
	"lifememo/application/moment/rpc/internal/svc"
	"lifememo/application/moment/rpc/pb"
)

type MomentServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedMomentServer
}

func NewMomentServer(svcCtx *svc.ServiceContext) *MomentServer {
	return &MomentServer{
		svcCtx: svcCtx,
	}
}

func (s *MomentServer) Publish(ctx context.Context, in *pb.PublishRequest) (*pb.PublishResponse, error) {
	l := logic.NewPublishLogic(ctx, s.svcCtx)
	return l.Publish(in)
}

func (s *MomentServer) Moments(ctx context.Context, in *pb.MomentsRequest) (*pb.MomentsResponse, error) {
	l := logic.NewMomentsLogic(ctx, s.svcCtx)
	return l.Moments(in)
}

func (s *MomentServer) MomentDelete(ctx context.Context, in *pb.MomentDeleteRequest) (*pb.MomentDeleteResponse, error) {
	l := logic.NewMomentDeleteLogic(ctx, s.svcCtx)
	return l.MomentDelete(in)
}

func (s *MomentServer) MomentDetail(ctx context.Context, in *pb.MomentDetailRequest) (*pb.MomentDetailResponse, error) {
	l := logic.NewMomentDetailLogic(ctx, s.svcCtx)
	return l.MomentDetail(in)
}
