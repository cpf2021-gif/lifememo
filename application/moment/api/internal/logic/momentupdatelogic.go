package logic

import (
	"context"
	"encoding/json"
	"lifememo/application/moment/api/internal/code"
	"lifememo/application/moment/rpc/pb"
	"lifememo/pkg/xcode"

	"lifememo/application/moment/api/internal/svc"
	"lifememo/application/moment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MomentUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMomentUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MomentUpdateLogic {
	return &MomentUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MomentUpdateLogic) MomentUpdate(req *types.MomentUpdateRequest) (*types.MomentUpdateResponse, error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("l.ctx.Value.error: %v", err)
		return nil, xcode.NoLogin
	}

	findMomentResp, err := l.svcCtx.MomentRPC.MomentDetail(l.ctx, &pb.MomentDetailRequest{
		MomentId: req.MomentId,
	})

	if err != nil {
		return nil, code.MomentNotFound
	}

	if findMomentResp.Moment.AuthorId != userId {
		return nil, code.PermissionDenied
	}

	_, err = l.svcCtx.MomentRPC.MomentUpdateContent(l.ctx, &pb.MomentUpdateContentRequest{
		MomentId: req.MomentId,
		Content:  req.Content,
	})

	if err != nil {
		return nil, err
	}

	return &types.MomentUpdateResponse{}, nil
}
