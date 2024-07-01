package logic

import (
	"context"
	"encoding/json"
	"lifememo/application/follow/api/internal/code"
	"lifememo/application/follow/rpc/follow"
	"lifememo/application/user/rpc/user"
	"lifememo/pkg/xcode"

	"lifememo/application/follow/api/internal/svc"
	"lifememo/application/follow/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnfollowUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnfollowUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnfollowUserLogic {
	return &UnfollowUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnfollowUserLogic) UnfollowUser(req *types.UnfollowRequest) (resp *types.UnfollowResponse, err error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("l.ctx.Value.error: %v", err)
		return nil, xcode.NoLogin
	}

	if req.FollowedUserId <= 0 {
		return nil, code.FollowUserIdInvalid
	}

	_, err = l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{
		UserId: req.FollowedUserId,
	})

	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.FollowRPC.UnFollow(l.ctx, &follow.UnFollowRequest{
		UserId:         userId,
		FollowedUserId: req.FollowedUserId,
	})

	if err != nil {
		return nil, err
	}

	return &types.UnfollowResponse{}, nil
}
