package logic

import (
	"context"
	"encoding/json"

	"lifememo/application/moment/api/internal/svc"
	"lifememo/application/moment/api/internal/types"
	"lifememo/application/moment/rpc/pb"
	"lifememo/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishLogic) Publish(req *types.PublishRequest) (*types.PublishResponse, error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("l.ctx.Value.error: %v", err)
		return nil, xcode.NoLogin
	}

	res, err := l.svcCtx.MomentRPC.Publish(l.ctx, &pb.PublishRequest{
		UserId:  userId,
		Content: req.Content,
	})
	if err != nil {
		logx.Error("l.svcCtx.MomentRPC.Publish req: %v userid: %d error: %v", req, userId, err)
		return nil, err
	}

	return &types.PublishResponse{
		MonmentId: res.MomentId,
	}, nil
}
