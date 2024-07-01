package logic

import (
	"context"
	"lifememo/application/moment/rpc/pb"
	"lifememo/application/user/rpc/user"

	"lifememo/application/moment/api/internal/svc"
	"lifememo/application/moment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MomentDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMomentDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MomentDetailLogic {
	return &MomentDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MomentDetailLogic) MomentDetail(req *types.MomentDetailRequest) (*types.MomentDetailResponse, error) {
	resp, err := l.svcCtx.MomentRPC.MomentDetail(l.ctx, &pb.MomentDetailRequest{
		MomentId: req.MomentId,
	})

	if err != nil {
		return nil, err
	}

	var moment *types.MomentItem
	moment = &types.MomentItem{
		Id:           resp.Moment.Id,
		AuthorId:     resp.Moment.AuthorId,
		Content:      resp.Moment.Content,
		CommentCount: resp.Moment.CommentCount,
		LikeCount:    resp.Moment.LikeCount,
		PublishTime:  resp.Moment.PublishTime,
	}

	User, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{
		UserId: resp.Moment.AuthorId,
	})
	if err != nil {
		return nil, err
	}

	moment.AuthorName = User.Username

	return &types.MomentDetailResponse{
		Moment: moment,
	}, nil
}
