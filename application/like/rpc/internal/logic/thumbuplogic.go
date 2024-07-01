package logic

import (
	"context"
	"encoding/json"
	"errors"
	"lifememo/application/like/rpc/internal/model"

	"lifememo/application/like/rpc/internal/svc"
	"lifememo/application/like/rpc/internal/types"
	"lifememo/application/like/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/threading"
)

type ThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

const (
	LikeCronKey = "likeCron"
)

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ThumbupLogic) Thumbup(in *service.ThumbupRequest) (*service.ThumbupResponse, error) {
	likeType := int64(in.LikeType)
	if likeType != types.Like && likeType != types.UnLike {
		return &service.ThumbupResponse{}, errors.New("likeType is invalid")
	}

	msg := &types.ThumbupMsg{
		BizId:    in.BizId,
		ObjId:    in.ObjId,
		UserId:   in.UserId,
		LikeType: int64(in.LikeType),
	}

	err := l.UpdateLikeNum(l.ctx, msg)
	if err != nil {
		l.Logger.Errorf("UpdateLikeNum err: %v", err)
		return &service.ThumbupResponse{}, err
	}

	// 异步更新点赞数
	threading.GoSafe(func() {
		data, err := json.Marshal(msg)
		if err != nil {
			logx.Errorf("[Thumbup] json.Marshal err: %v", err)
			return
		}
		if msg.BizId == "moment" {
			err = l.svcCtx.MomentLikeKqPusherClient.Push(string(data))
			if err != nil {
				logx.Errorf("[Thumbup] Push err: %v", err)
				return
			}
		} else {
			err = l.svcCtx.ReplyLikeKqPusherClient.Push(string(data))
			if err != nil {
				logx.Errorf("[Thumbup] Push err: %v", err)
				return
			}
		}
	})

	return &service.ThumbupResponse{}, nil
}

func (l *ThumbupLogic) UpdateLikeNum(ctx context.Context, msg *types.ThumbupMsg) error {
	// 事务
	return l.svcCtx.LikeRecordModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 保存点赞记录
		// 先查看记录是否存在
		recordRes, err := l.svcCtx.LikeRecordModel.FindOneByBizIdObjIdUserId(ctx, msg.BizId, msg.ObjId, msg.UserId)
		if err == model.ErrNotFound {
			// 不存在，直接插入
			// 如果是取消点赞，不需要插入
			if msg.LikeType == types.UnLike {
				return nil
			} else {
				// 如果是点赞，需要插入
				if _, err = l.svcCtx.LikeRecordModel.InsertWithSession(ctx, session, &model.LikeRecord{
					BizId:    msg.BizId,
					ObjId:    msg.ObjId,
					UserId:   msg.UserId,
					LikeType: msg.LikeType,
				}); err != nil {
					return err
				}
			}
		} else if err != nil {
			return err
		} else {
			// 存在，判断是否需要更新
			if recordRes.LikeType == msg.LikeType { // likeType不变，不需要更新
				return nil
			} else { // likeType变更，需要更新
				if err = l.svcCtx.LikeRecordModel.UpdateWithSession(ctx, session, &model.LikeRecord{
					Id:       recordRes.Id,
					BizId:    msg.BizId,
					ObjId:    msg.ObjId,
					UserId:   msg.UserId,
					LikeType: msg.LikeType,
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
			if msg.LikeType == types.UnLike { // 取消点赞，需要减1
				if likeCountRes.LikeNum == 0 {
					return nil
				}
				if err = l.svcCtx.LikeCountModel.UpdateWithSession(ctx, session, &model.LikeCount{
					Id:      likeCountRes.Id,
					BizId:   msg.BizId,
					ObjId:   msg.ObjId,
					LikeNum: likeCountRes.LikeNum - 1,
				}); err != nil {
					return err
				}
			} else { // 点赞，需要加1
				if err = l.svcCtx.LikeCountModel.UpdateWithSession(ctx, session, &model.LikeCount{
					Id:      likeCountRes.Id,
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
