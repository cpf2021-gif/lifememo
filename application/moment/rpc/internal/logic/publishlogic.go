package logic

import (
	"context"
	"time"

	"lifememo/application/moment/rpc/internal/code"
	"lifememo/application/moment/rpc/internal/model"
	"lifememo/application/moment/rpc/internal/svc"
	"lifememo/application/moment/rpc/internal/types"
	"lifememo/application/moment/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishLogic) Publish(in *pb.PublishRequest) (*pb.PublishResponse, error) {
	if in.UserId <= 0 {
		return nil, code.MomentIdInvalid
	}

	if len(in.Content) == 0 {
		return nil, code.MomentContentEmpty
	}

	ret, err := l.svcCtx.MomentModel.Insert(l.ctx, &model.Moment{
		AuthorId:    in.UserId,
		Content:     in.Content,
		Status:      types.MomentStatusVisible,
		PublishTime: time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("Publish Insert req: %v error: %v", in, err)
		return nil, err
	}

	momentId, err := ret.LastInsertId()
	if err != nil {
		l.Logger.Errorf("LastInsertId error: %v", err)
		return nil, err
	}

	return &pb.PublishResponse{MomentId: momentId}, nil
}
