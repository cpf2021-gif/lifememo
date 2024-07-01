package logic

import (
	"context"
	"encoding/json"
	"lifememo/application/message/api/internal/code"
	"lifememo/application/message/rpc/messageservice"
	"lifememo/application/user/rpc/user"
	"lifememo/pkg/xcode"

	"lifememo/application/message/api/internal/svc"
	"lifememo/application/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SentMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSentMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SentMessageLogic {
	return &SentMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SentMessageLogic) SentMessage(req *types.SentMessageRequest) (*types.SentMessageResponse, error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("l.ctx.Value.error: %v", err)
		return nil, xcode.NoLogin
	}

	if userId == req.ReceiverId {
		return nil, code.NoSendToSelf
	}

	// 检测对象是否存在
	_, err = l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{
		UserId: req.ReceiverId,
	})

	if err != nil {
		return nil, err
	}

	// 发送消息
	resp, err := l.svcCtx.MessageRPC.SentMessage(l.ctx, &messageservice.SentMessageRequest{
		SenderId:   userId,
		ReceiverId: req.ReceiverId,
		Content:    req.Content,
	})

	if err != nil {
		return nil, err
	}

	return &types.SentMessageResponse{
		MessageId: resp.MessageId,
	}, nil
}
