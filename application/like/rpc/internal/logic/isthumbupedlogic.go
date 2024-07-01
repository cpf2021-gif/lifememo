package logic

import (
	"context"
	"errors"
	"lifememo/application/like/rpc/internal/model"
	"lifememo/application/like/rpc/internal/types"

	"lifememo/application/like/rpc/internal/svc"
	"lifememo/application/like/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsThumbupedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsThumbupedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsThumbupedLogic {
	return &IsThumbupedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsThumbupedLogic) IsThumbuped(in *service.IsThumbupedRequest) (*service.IsThumbupedResponse, error) {
	reply, err := l.svcCtx.LikeRecordModel.FindOneByBizIdObjIdUserId(l.ctx, in.BizId, in.ObjId, in.UserId)
	if errors.Is(err, model.ErrNotFound) {
		return &service.IsThumbupedResponse{IsThumbuped: false}, nil
	} else if err != nil {
		return &service.IsThumbupedResponse{}, err
	} else {
		return &service.IsThumbupedResponse{IsThumbuped: reply.LikeType == types.Like}, nil
	}
}
