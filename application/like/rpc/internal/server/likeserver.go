// Code generated by goctl. DO NOT EDIT.
// Source: like.proto

package server

import (
	"context"

	"lifememo/application/like/rpc/internal/logic"
	"lifememo/application/like/rpc/internal/svc"
	"lifememo/application/like/rpc/service"
)

type LikeServer struct {
	svcCtx *svc.ServiceContext
	service.UnimplementedLikeServer
}

func NewLikeServer(svcCtx *svc.ServiceContext) *LikeServer {
	return &LikeServer{
		svcCtx: svcCtx,
	}
}

func (s *LikeServer) Thumbup(ctx context.Context, in *service.ThumbupRequest) (*service.ThumbupResponse, error) {
	l := logic.NewThumbupLogic(ctx, s.svcCtx)
	return l.Thumbup(in)
}