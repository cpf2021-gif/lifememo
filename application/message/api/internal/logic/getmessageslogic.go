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

type GetMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessagesLogic {
	return &GetMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMessagesLogic) GetMessages(req *types.GetMessagesRequest) (*types.GetMessagesResponse, error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("l.ctx.Value.error: %v", err)
		return nil, xcode.NoLogin
	}

	if userId == req.TargetId {
		return nil, code.NoSendToSelf
	}

	// 检测对象是否存在
	_, err = l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{
		UserId: req.TargetId,
	})

	if err != nil {
		return nil, err
	}

	// 获取消息
	resp, err := l.svcCtx.MessageRPC.GetMessages(l.ctx, &messageservice.GetMessagesRequest{
		UserId:   userId,
		TargetId: req.TargetId,
	})

	if err != nil {
		return nil, err
	}

	var messages []*types.MessageItem
	for _, v := range resp.Messages {
		messages = append(messages, &types.MessageItem{
			Id:          v.MessageId,
			SenderId:    v.SenderId,
			ReceiverId:  v.ReceiverId,
			Content:     v.Content,
			CreatedTime: v.CreatedTime,
		})
	}

	return &types.GetMessagesResponse{
		Messages: messages,
	}, nil
}
