package logic

import (
	"context"
	"encoding/json"
	"lifememo/application/message/rpc/messageservice"
	"lifememo/pkg/xcode"

	"lifememo/application/message/api/internal/svc"
	"lifememo/application/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMessageLogic {
	return &DeleteMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMessageLogic) DeleteMessage(req *types.DeleteMessageRequest) (*types.DeleteMessageResponse, error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("l.ctx.Value.error: %v", err)
		return nil, xcode.NoLogin
	}

	_, err = l.svcCtx.MessageRPC.DeleteMessage(l.ctx, &messageservice.DeleteMessageRequest{
		MessageId: req.MessageId,
		UserId:    userId,
	})

	if err != nil {
		return nil, err
	}

	return &types.DeleteMessageResponse{}, nil
}
