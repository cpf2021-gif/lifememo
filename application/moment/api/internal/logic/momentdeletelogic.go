package logic

import (
	"context"
	"encoding/json"
	"lifememo/application/moment/rpc/pb"
	"lifememo/pkg/xcode"

	"lifememo/application/moment/api/internal/svc"
	"lifememo/application/moment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MomentDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMomentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MomentDeleteLogic {
	return &MomentDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MomentDeleteLogic) MomentDelete(req *types.MomentDeleteRequest) (*types.MomentDeleteResponse, error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("l.ctx.Value.error: %v", err)
		return nil, xcode.NoLogin
	}

	_, err = l.svcCtx.MomentRPC.MomentDelete(l.ctx, &pb.MomentDeleteRequest{
		MomentId: req.MomentId,
		UserId:   userId,
	})

	if err != nil {
		return nil, err
	}

	return &types.MomentDeleteResponse{}, nil
}
