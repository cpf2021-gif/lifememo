package logic

import (
	"context"
	"encoding/json"
	"lifememo/application/like/cron/internal/model"
	"lifememo/application/like/cron/internal/svc"
	"lifememo/application/like/cron/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
	pairs, err := l.svcCtx.BizRedis.ZrangebyscoreWithScoresCtx(l.ctx, LikeCronKey, 0, 2)
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
			continue
		}

		err = l.UpdateLikeNum(l.ctx, msg, v)
		if err != nil {
			l.Logger.Errorf("UpdateLikeNum error:%v", err)
			continue
		}
	}

	_, err = l.svcCtx.BizRedis.Del(LikeCronKey)
	if err != nil {
		l.Logger.Errorf("Del error:%v", err)
	}
}

func (l *UpdateLogic) UpdateLikeNum(ctx context.Context, msg *types.ThumbupMsg, v int64) error {
	return l.svcCtx.LikeRecordModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 保存点赞记录
		// 先查看记录是否存在
		recordRes, err := l.svcCtx.LikeRecordModel.FindOneByBizIdObjIdUserId(ctx, msg.BizId, msg.ObjId, msg.UserId)
		if err == model.ErrNotFound {
			// 不存在，直接插入
			// 如果是取消点赞，不需要插入
			if v == types.UnLike {
				return nil
			} else {
				if _, err = l.svcCtx.LikeRecordModel.InsertWithSession(ctx, session, &model.LikeRecord{
					BizId:    msg.BizId,
					ObjId:    msg.ObjId,
					UserId:   msg.UserId,
					LikeType: v,
				}); err != nil {
					return err
				}

			}
		} else if err != nil {
			return err
		} else {
			// 存在，判断是否需要更新
			if recordRes.LikeType == v { // likeType不变，不需要更新
				return nil
			} else { // likeType变更，需要更新
				if err = l.svcCtx.LikeRecordModel.UpdateWithSession(ctx, session, &model.LikeRecord{
					Id:       recordRes.Id,
					BizId:    msg.BizId,
					ObjId:    msg.ObjId,
					UserId:   msg.UserId,
					LikeType: v,
				}); err != nil {
					return err
				}
			}
		}

		// 保存点赞数
		// 查看动态点赞数是否存在
		likeCountRes, err := l.svcCtx.LikeCountModel.FindOneByBizIdObjId(ctx, msg.BizId, msg.ObjId)
		if err == model.ErrNotFound {
			// 不存在，直接插入
			if _, err = l.svcCtx.LikeCountModel.InsertWithSession(ctx, session, &model.LikeCount{
				BizId:   msg.BizId,
				ObjId:   msg.ObjId,
				LikeNum: 1,
			}); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			// 存在，判断是否需要更新
			if v == types.UnLike { // 取消点赞，需要减1
				if likeCountRes.LikeNum == 0 {
					return nil
				}
				if err = l.svcCtx.LikeCountModel.UpdateWithSession(ctx, session, &model.LikeCount{
					// Id:      likeCountRes.Id,
					BizId:   msg.BizId,
					ObjId:   msg.ObjId,
					LikeNum: likeCountRes.LikeNum - 1,
				}); err != nil {
					return err
				}
			} else { // 点赞，需要加1
				if err = l.svcCtx.LikeCountModel.UpdateWithSession(ctx, session, &model.LikeCount{
					// Id:      likeCountRes.Id,
					BizId:   msg.BizId,
					ObjId:   msg.ObjId,
					LikeNum: likeCountRes.LikeNum + 1,
				}); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
