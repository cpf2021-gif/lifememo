package logic

import (
	"context"
	"errors"
	"lifememo/application/moment/rpc/internal/model"

	"lifememo/application/moment/rpc/internal/svc"
	"lifememo/application/moment/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MomentUpdateContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMomentUpdateContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MomentUpdateContentLogic {
	return &MomentUpdateContentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MomentUpdateContentLogic) MomentUpdateContent(in *pb.MomentUpdateContentRequest) (*pb.MomentUpdateContentResponse, error) {
	moment, err := l.svcCtx.MomentModel.FindOne(l.ctx, in.MomentId)
	if errors.Is(err, model.ErrNotFound) {
		return &pb.MomentUpdateContentResponse{}, errors.New("动态不存在")
	} else if err != nil {
		return &pb.MomentUpdateContentResponse{}, err
	} else {
		moment.Content = in.Content
		err = l.svcCtx.MomentModel.Update(l.ctx, moment)
		if err != nil {
			return &pb.MomentUpdateContentResponse{}, err
		}
	}
	return &pb.MomentUpdateContentResponse{}, nil
}
