package logic

import (
	"context"
	"encoding/json"
	"lifememo/application/like/mq/internal/svc"
	"lifememo/application/like/mq/internal/types"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

type MomentLikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMomentLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MomentLikeLogic {
	return &MomentLikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MomentLikeLogic) Consume(_, val string) error {
	logx.Infof("[MomentLike] Consume: %s", val)
	var msg *types.ThumbupMsg
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		logx.Errorf("[MomentLike] Consume val: %s, err: %v", val, err)
		return err
	}

	moment, err := l.svcCtx.MomentModel.FindOne(context.Background(), msg.ObjId)
	if err != nil {
		logx.Errorf("[MomentLike] Consume FindOne err: %v", err)
		return err
	}

	if msg.LikeType == types.Like {
		moment.LikeNum++
	} else {
		moment.LikeNum--
		if moment.LikeNum < 0 {
			moment.LikeNum = 0
			logx.Errorf("[MomentLike] Consume moment.LikeNum < 0, moment: %v", moment)
		}
	}

	err = l.svcCtx.MomentModel.Update(context.Background(), moment)
	if err != nil {
		logx.Errorf("[MomentLike] Consume Update err: %v", err)
	}

	return err
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.MomentLikeKqConsumerConf, NewMomentLikeLogic(ctx, svcCtx)),
		kq.MustNewQueue(svcCtx.Config.ReplyLikeKqConsumerConf, NewReplyLikeLogic(ctx, svcCtx)),
	}
}
