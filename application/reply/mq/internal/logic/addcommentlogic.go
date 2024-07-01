package logic

import (
	"context"
	"lifememo/application/reply/mq/internal/svc"
	"strconv"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

type AddCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCommentLogic {
	return &AddCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddCommentLogic) Consume(_, val string) error {
	logx.Infof("[AddComment] Consume: %s", val)
	momentId, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		logx.Errorf("[AddComment] strconv.ParseInt err: %v", err)
	}
	moment, err := l.svcCtx.MomentModel.FindOne(l.ctx, momentId)
	if err != nil {
		logx.Errorf("[AddComment] MomentModel.FindOne err: %v", err)
	}
	moment.CommentNum++
	err = l.svcCtx.MomentModel.Update(l.ctx, moment)
	if err != nil {
		logx.Errorf("[AddComment] MomentModel.Update err: %v", err)
	}
	return nil
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.MomentReplyKqConsumerConf, NewAddCommentLogic(ctx, svcCtx)),
	}
}
