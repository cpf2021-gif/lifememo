package logic

import (
	"context"
	"encoding/json"

	"lifememo/application/applet/api/internal/svc"
	"lifememo/application/applet/api/internal/types"
	"lifememo/application/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (*types.UserInfoResponse, error) {
	userId, err := l.ctx.Value(types.UserIdKey).(json.Number).Int64()
	if err != nil {
		return nil, err
	}
	if userId == 0 {
		return nil, nil
	}
	u, err := l.svcCtx.UserRpc.FindById(l.ctx, &user.FindByIdRequest{
		UserId: userId,
	})
	if err != nil {
		logx.Errorf("FindById error: %v", userId)
		return nil, err
	}

	return &types.UserInfoResponse{
		UserId:   u.UserId,
		Username: u.Username,
	}, nil
}
