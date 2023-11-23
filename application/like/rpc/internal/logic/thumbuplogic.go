package logic

import (
	"context"
	"encoding/json"

	"lifememo/application/like/rpc/internal/svc"
	"lifememo/application/like/rpc/internal/types"
	"lifememo/application/like/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type ThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

const (
	LikeCronKey = "likeCron"
)

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ThumbupLogic) Thumbup(in *service.ThumbupRequest) (*service.ThumbupResponse, error) {
	msg := &types.ThumbupMsg{
		BizId:  in.BizId,
		ObjId:  in.ObjId,
		UserId: in.UserId,
	}

	s, err := json.Marshal(msg)
	if err != nil {
		l.Logger.Errorf("json Marshal val:%v err: %v", msg, err)
		return &service.ThumbupResponse{}, nil
	}
	_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, LikeCronKey, int64(in.LikeType), string(s))
	if err != nil {
		l.Logger.Errorf("ZaddCtx err: %v", err)
	}

	return &service.ThumbupResponse{}, nil
}
