package logic

import (
	"context"
	"lifememo/application/moment/rpc/internal/code"
	"lifememo/application/moment/rpc/internal/model"
	"lifememo/application/moment/rpc/internal/svc"
	"lifememo/application/moment/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MomentDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMomentDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MomentDetailLogic {
	return &MomentDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MomentDetailLogic) MomentDetail(in *pb.MomentDetailRequest) (*pb.MomentDetailResponse, error) {
	moment, err := l.svcCtx.MomentModel.FindOne(l.ctx, in.MomentId)
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.MomentDetailResponse{}, code.MomentNotFound
		}
		return nil, err
	}

	if err = l.svcCtx.MomentModel.AddMomentViewNum(l.ctx, in.MomentId); err != nil {
		l.Logger.Errorf("AddMomentViewNum error: %v", err)
	}

	return &pb.MomentDetailResponse{
		Moment: &pb.MomentItem{
			Id:           moment.Id,
			Content:      moment.Content,
			CommentCount: moment.CommentNum,
			LikeCount:    moment.LikeNum,
			PublishTime:  moment.PublishTime.Unix(),
			AuthorId:     moment.AuthorId,
		},
	}, nil
}
