package logic

import (
	"context"
	"lifememo/application/follow/rpc/follow"
	"lifememo/application/user/rpc/user"

	"lifememo/application/follow/api/internal/svc"
	"lifememo/application/follow/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowListLogic {
	return &GetFollowListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFollowListLogic) GetFollowList(req *types.GetFollowListRequest) (*types.GetFollowListResponse, error) {
	_, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	resp, err := l.svcCtx.FollowRPC.FollowList(l.ctx, &follow.FollowListRequest{
		Id:       req.Id,
		UserId:   req.UserId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})

	var followList []*types.FollowItem
	for _, item := range resp.Items {
		followList = append(followList, &types.FollowItem{
			FollowedUserId: item.FollowedUserId,
			FansCount:      item.FansCount,
			CreateTime:     item.CreateTime,
		})
	}

	return &types.GetFollowListResponse{
		FollowList: followList,
		Cursor:     resp.Cursor,
		IsEnd:      resp.IsEnd,
		Id:         resp.Id,
	}, nil
}
