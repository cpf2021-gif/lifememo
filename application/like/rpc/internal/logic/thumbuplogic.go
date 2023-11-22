package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"lifememo/application/like/rpc/internal/types"
	"strconv"

	"lifememo/application/like/rpc/internal/svc"
	"lifememo/application/like/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type ThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

const (
	prefixLikeRecord = "biz#like#%s#%d#%d"
)

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ThumbupLogic) Thumbup(in *service.ThumbupRequest) (*service.ThumbupResponse, error) {
	msg := &types.ThumbupMsg{
		BizId:    in.BizId,
		ObjId:    in.ObjId,
		UserId:   in.UserId,
		LikeType: in.LikeType,
	}

	key := likeRecordKey(in.BizId, in.ObjId, in.UserId)

	var likeType int64 = 2
	isUpdate := false

	// 是否存在
	b, err := l.svcCtx.BizRedis.ExistsCtx(l.ctx, key)
	if err != nil {
		l.Logger.Errorf("Redis ExistsCtx val:%s error:%v", key, err)
	}

	if b {
		val, err := l.svcCtx.BizRedis.GetCtx(l.ctx, key)
		if err != nil {
			l.Logger.Errorf("Redis GetCtx val:%s error:%v", key, err)
		}
		if len(val) != 0 {
			res, err := strconv.Atoi(val)
			if err != nil {
				l.Logger.Errorf("Atoi val:%s error:%v", val, err)
			}
			likeType = int64(res)
		}
		if likeType != 2 && likeType != int64(msg.LikeType) {
			isUpdate = true
			val := strconv.Itoa(int(msg.LikeType))
			err := l.svcCtx.BizRedis.SetCtx(l.ctx, key, val)
			if err != nil {
				l.Logger.Errorf("Redis SetCtx key:%s val:%s error:%v", key, val, err)
			}
		}
	} else {
		val := strconv.Itoa(int(msg.LikeType))
		err := l.svcCtx.BizRedis.SetCtx(l.ctx, key, val)
		if err != nil {
			l.Logger.Errorf("Redis SetCtx key:%s val:%s error:%v", key, val, err)
		}
	}

	if !b || isUpdate {
		threading.GoSafe(func() {
			data, err := json.Marshal(msg)
			if err != nil {
				l.Logger.Errorf("[Thumbup] marshal msg: %v error: %v", msg, err)
				return
			}
			err = l.svcCtx.KqPusherClient.Push(string(data))
			if err != nil {
				l.Logger.Errorf("[Thumbup] kq push data: %s error: %v", data, err)
			}
		})
	}

	return &service.ThumbupResponse{}, nil
}

func likeRecordKey(bizId string, objId, userId int64) string {
	return fmt.Sprintf(prefixLikeRecord, bizId, objId, userId)
}
