package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"lifememo/application/like/mq/internal/model"
	"lifememo/application/like/mq/internal/svc"
	"lifememo/application/like/mq/internal/types"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

type ThumupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumupLogic {
	return &ThumupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ThumupLogic) Consume(key, val string) error {
	fmt.Printf("get key: %s val: %s\n", key, val)

	// string -> struct
	jsonString := val
	var likeMsg *types.ThumbupMsg

	err := json.Unmarshal([]byte(jsonString), &likeMsg)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return nil
	}

	return l.updateLikeNum(l.ctx, likeMsg)
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, NewThumbupLogic(ctx, svcCtx)),
	}
}

func (l *ThumupLogic) updateLikeNum(ctx context.Context, msg *types.ThumbupMsg) error {
	res, err := l.svcCtx.LikeRecordModel.FindOneByBizIdObjIdUserId(ctx, msg.BizId, msg.ObjId, msg.UserId)
	if err == model.ErrNotFound {
		_, err = l.svcCtx.LikeRecordModel.Insert(ctx, &model.LikeRecord{
			BizId:    msg.BizId,
			ObjId:    msg.ObjId,
			UserId:   msg.UserId,
			LikeType: int64(msg.LikeType),
		})
		return err
	}
	if err == nil {
		if res.LikeType == int64(msg.LikeType) {
			return nil
		}
		res.UpdateTime = time.Now()
		res.LikeType = int64(msg.LikeType)
		err = l.svcCtx.LikeRecordModel.Update(ctx, res)
	}
	return err
}
