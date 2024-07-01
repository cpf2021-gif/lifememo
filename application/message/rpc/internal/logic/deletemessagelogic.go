package logic

import (
	"context"
	"lifememo/application/message/rpc/internal/code"
	"lifememo/application/message/rpc/internal/model"
	"lifememo/application/message/rpc/internal/types"

	"lifememo/application/message/rpc/internal/svc"
	"lifememo/application/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMessageLogic {
	return &DeleteMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMessageLogic) DeleteMessage(in *pb.DeleteMessageRequest) (*pb.DeleteMessageResponse, error) {
	message, err := l.svcCtx.MessageModel.FindOne(l.ctx, in.MessageId)
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.DeleteMessageResponse{}, code.MessageInvalid
		}
		return &pb.DeleteMessageResponse{}, err
	}

	// 只能删除自己的消息
	if message.SendUserId != in.UserId {
		return &pb.DeleteMessageResponse{}, code.NoPermission
	}

	// 已经删除的消息不再删除
	if message.Status == types.DeletedStatus {
		return &pb.DeleteMessageResponse{}, code.DeletedMessage
	} else {
		message.Status = types.DeletedStatus
		err = l.svcCtx.MessageModel.Update(l.ctx, message)
		if err != nil {
			return &pb.DeleteMessageResponse{}, err
		}
	}

	return &pb.DeleteMessageResponse{}, nil
}
