package logic

import (
	"context"
	"lifememo/application/message/rpc/internal/code"
	"lifememo/application/message/rpc/internal/model"

	"lifememo/application/message/rpc/internal/svc"
	"lifememo/application/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessagesLogic {
	return &GetMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessagesLogic) GetMessages(in *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	messages, err := l.svcCtx.MessageModel.FindMessageBySenderAndReceiver(in.UserId, in.TargetId)
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.GetMessagesResponse{}, code.NoneMessage
		}
		return &pb.GetMessagesResponse{}, err
	}

	var res []*pb.MessageItem
	for _, message := range messages {
		res = append(res, &pb.MessageItem{
			MessageId:   message.Id,
			SenderId:    message.SendUserId,
			ReceiverId:  message.ReceiveUserId,
			Content:     message.Content,
			CreatedTime: message.CreateTime.Unix(),
		})
	}

	return &pb.GetMessagesResponse{
		Messages: res,
	}, nil
}
