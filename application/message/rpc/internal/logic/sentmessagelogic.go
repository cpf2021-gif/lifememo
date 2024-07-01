package logic

import (
	"context"
	"lifememo/application/message/rpc/internal/model"
	"lifememo/application/message/rpc/internal/types"

	"lifememo/application/message/rpc/internal/svc"
	"lifememo/application/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SentMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSentMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SentMessageLogic {
	return &SentMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SentMessageLogic) SentMessage(in *pb.SentMessageRequest) (*pb.SentMessageResponse, error) {
	res, err := l.svcCtx.MessageModel.Insert(l.ctx, &model.Message{
		SendUserId:    in.SenderId,
		ReceiveUserId: in.ReceiverId,
		Content:       in.Content,
		Status:        types.NormalStatus,
	})

	lastId, err := res.LastInsertId()
	if err != nil {
		return &pb.SentMessageResponse{}, err
	}

	if err != nil {
		return &pb.SentMessageResponse{}, err
	}

	return &pb.SentMessageResponse{
		MessageId: lastId,
	}, nil
}
