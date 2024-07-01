package logic

import (
	"context"
	"lifememo/application/follow/rpc/follow"
	"lifememo/application/user/rpc/user"

	"lifememo/application/follow/api/internal/svc"
	"lifememo/application/follow/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFansListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFansListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFansListLogic {
	return &GetFansListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFansListLogic) GetFansList(req *types.GetFansListRequest) (*types.GetFansListResponse, error) {
	_, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	resp, err := l.svcCtx.FollowRPC.FansList(l.ctx, &follow.FansListRequest{
		Id:       req.Id,
		UserId:   req.UserId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})

	var fansList []*types.FansItem
	for _, item := range resp.Items {
		fansList = append(fansList, &types.FansItem{
			FansUserId: item.FansUserId,
			FansCount:  item.FansCount,
			CreateTime: item.CreateTime,
		})
	}

	return &types.GetFansListResponse{
		FansList: fansList,
		Cursor:   resp.Cursor,
		IsEnd:    resp.IsEnd,
		Id:       resp.Id,
	}, nil
}
