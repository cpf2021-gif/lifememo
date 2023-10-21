package logic

import (
	"context"

	"lifememo/application/user/rpc/internal/svc"
	"lifememo/application/user/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByEmailLogic {
	return &FindByEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByEmailLogic) FindByEmail(in *service.FindByEmailRequest) (*service.FindByEmailResponse, error) {
	user, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, in.Email)
	if err != nil {
		logx.Errorf("FindByEmail error: %v", err)
	}

	if user == nil {
		return &service.FindByEmailResponse{}, nil
	}

	return &service.FindByEmailResponse{
		UserId:   user.Id,
		Username: user.Username,
	}, nil
}
