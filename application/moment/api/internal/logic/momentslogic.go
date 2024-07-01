package logic

import (
	"context"
	"lifememo/application/moment/api/internal/svc"
	"lifememo/application/moment/api/internal/types"
	"lifememo/application/moment/rpc/pb"
	"lifememo/application/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type MomentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMomentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MomentsLogic {
	return &MomentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MomentsLogic) Moments(req *types.MomentsRequest) (*types.MomentsResponse, error) {
	// 获取用户信息
	User, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	moments, err := l.svcCtx.MomentRPC.Moments(l.ctx, &pb.MomentsRequest{
		UserId:   req.UserId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
		SortType: req.SortType,
		MomentId: req.MomentId,
	})
	if err != nil {
		return nil, err
	}

	var momentList []*types.MomentItem
	for _, v := range moments.Moments {
		momentList = append(momentList, &types.MomentItem{
			Id:           v.Id,
			AuthorId:     v.AuthorId,
			Content:      v.Content,
			CommentCount: v.CommentCount,
			LikeCount:    v.LikeCount,
			PublishTime:  v.PublishTime,
		})

	}

	for _, v := range momentList {
		v.AuthorName = User.Username
	}

	return &types.MomentsResponse{
		Moments:  momentList,
		Cursor:   moments.Cursor,
		IsEnd:    moments.IsEnd,
		MomentId: moments.MomentId,
	}, nil
}
