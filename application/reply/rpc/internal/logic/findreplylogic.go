package logic

import (
	"context"
	"errors"
	"lifememo/application/reply/rpc/internal/code"
	"lifememo/application/reply/rpc/internal/model"

	"lifememo/application/reply/rpc/internal/svc"
	"lifememo/application/reply/rpc/pb1"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindReplyLogic {
	return &FindReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindReplyLogic) FindReply(in *pb1.FindReplyRequest) (*pb1.FindReplyResponse, error) {
	_, err := l.svcCtx.ReplyModel.FindOne(l.ctx, in.ReplyId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, code.ReplyNotFound
		}
		return nil, err
	}

	return &pb1.FindReplyResponse{}, nil
}
