package logic

import (
	"context"

	"lifememo/application/moment/rpc/internal/code"
	"lifememo/application/moment/rpc/internal/svc"
	"lifememo/application/moment/rpc/internal/types"
	"lifememo/application/moment/rpc/pb"
	"lifememo/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type MomentDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMomentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MomentDeleteLogic {
	return &MomentDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MomentDeleteLogic) MomentDelete(in *pb.MomentDeleteRequest) (*pb.MomentDeleteResponse, error) {
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}

	if in.MomentId <= 0 {
		return nil, code.MomentIdInvalid
	}

	moment, err := l.svcCtx.MomentModel.FindOne(l.ctx, in.MomentId)
	if err != nil {
		l.Logger.Errorf("MomentDelete FindOne req:%v, err:%v", in, err)
		return nil, err
	}
	if moment.AuthorId != in.UserId {
		return nil, xcode.AccessDenied
	}

	err = l.svcCtx.MomentModel.UpdateMomentStatus(l.ctx, in.MomentId, 4)
	if err != nil {
		l.Logger.Errorf("UpdateArticleStatus req:%v, err:%v", in, err)
		return nil, err
	}

	l.remCache(in, err)

	return &pb.MomentDeleteResponse{}, nil
}

func (l *MomentDeleteLogic) remCache(in *pb.MomentDeleteRequest, err error) {
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, momentsKey(in.UserId, types.SortPublishTime), in.MomentId)
	if err != nil {
		l.Logger.Errorf("ZremCtx req:%v, err:%v", in, err)
	}
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, momentsKey(in.UserId, types.SortPublishTime), in.MomentId)
	if err != nil {
		l.Logger.Errorf("ZremCtx req:%v, err:%v", in, err)
	}
}
