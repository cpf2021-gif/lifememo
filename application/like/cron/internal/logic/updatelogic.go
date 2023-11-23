package logic

import (
	"context"
	"encoding/json"
	"lifememo/application/like/cron/internal/model"
	"lifememo/application/like/cron/internal/svc"
	"lifememo/application/like/cron/internal/types"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

const (
	LikeCronKey = "likeCron"
)

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLogic) Update() {
	pairs, err := l.svcCtx.BizRedis.ZrangebyscoreWithScoresCtx(l.ctx, LikeCronKey, 0, 3)
	if err != nil {
		l.Logger.Errorf("ZrangebyscoreWithScoresCtx error:%v", err)
		return
	}

	for _, p := range pairs {
		k, v := p.Key, p.Score
		var msg *types.ThumbupMsg
		err = json.Unmarshal([]byte(k), &msg)
		if err != nil {
			l.Logger.Errorf("json.Unmarshal error:%v", err)
		}

		l.UpdateLikeNum(l.ctx, msg, v)
	}

	_, err = l.svcCtx.BizRedis.Del("likeCron")
	if err != nil {
		l.Logger.Errorf("Del error:%v", err)
	}
}

func (l *UpdateLogic) UpdateLikeNum(ctx context.Context, msg *types.ThumbupMsg, v int64) {

	res, err := l.svcCtx.LikeRecordModel.FindOneByBizIdObjIdUserId(ctx, msg.BizId, msg.ObjId, msg.UserId)
	if err == model.ErrNotFound {
		_, err = l.svcCtx.LikeRecordModel.Insert(ctx, &model.LikeRecord{
			BizId:      msg.BizId,
			ObjId:      msg.ObjId,
			UserId:     msg.UserId,
			LikeType:   v,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		})
		if err != nil {
			l.Logger.Errorf("Insert error:%v", err)
		}
		return
	}
	if err != nil {
		l.Logger.Errorf("FindOneByBizIdObjIdUserId error:%v", err)
		return
	}
	if res.LikeType != v {
		res.UpdateTime = time.Now()
		res.LikeType = v
		err = l.svcCtx.LikeRecordModel.Update(ctx, res)
		if err != nil {
			l.Logger.Errorf("Update error:%v", err)
		}
	}
}
