package logic

import (
	"context"

	"lifememo/application/user/rpc/internal/svc"
	"lifememo/application/user/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByIdLogic {
	return &FindByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByIdLogic) FindById(in *service.FindByIdRequest) (*service.FindByIdResponse, error) {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		logx.Errorf("FindById error: %v", err)
		return nil, err
	}

	return &service.FindByIdResponse{
		UserId:   user.Id,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
