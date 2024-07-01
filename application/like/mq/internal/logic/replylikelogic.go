package logic

import (
	"context"
	"encoding/json"
	"lifememo/application/like/mq/internal/svc"
	"lifememo/application/like/mq/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyLikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyLikeLogic {
	return &ReplyLikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReplyLikeLogic) Consume(_, val string) error {
	logx.Infof("[ReplyLike] Consume: %s", val)
	var msg *types.ThumbupMsg
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		logx.Errorf("[ReplyLike] Consume val: %s, err: %v", val, err)
		return err
	}

	reply, err := l.svcCtx.ReplyModel.FindOne(context.Background(), msg.ObjId)
	if err != nil {
		logx.Errorf("[ReplyLike] Consume FindOne err: %v", err)
		return err
	}

	if msg.LikeType == types.Like {
		reply.LikeNum++
	} else {
		reply.LikeNum--
		if reply.LikeNum < 0 {
			reply.LikeNum = 0
			logx.Errorf("[ReplyLike] Consume reply.LikeNum < 0, reply: %v", reply)
		}
	}

	err = l.svcCtx.ReplyModel.Update(context.Background(), reply)
	if err != nil {
		logx.Errorf("[ReplyLike] Consume Update err: %v", err)
	}
	return err
}
