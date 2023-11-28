package logic

import (
	"context"
	"strconv"
	"time"

	"lifememo/application/moment/rpc/internal/code"
	"lifememo/application/moment/rpc/internal/model"
	"lifememo/application/moment/rpc/internal/svc"
	"lifememo/application/moment/rpc/internal/types"
	"lifememo/application/moment/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishLogic) Publish(in *pb.PublishRequest) (*pb.PublishResponse, error) {
	if in.UserId <= 0 {
		return nil, code.MomentIdInvalid
	}

	if len(in.Content) == 0 {
		return nil, code.MomentContentEmpty
	}

	if in.Status != types.MomentStatusVisible && in.Status != types.MomentStatusPrivate {
		return nil, code.MomentStatusInvalid
	}

	now := time.Now()

	ret, err := l.svcCtx.MomentModel.Insert(l.ctx, &model.Moment{
		AuthorId:    in.UserId,
		Content:     in.Content,
		Status:      in.Status,
		PublishTime: now,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("Publish Insert req: %v error: %v", in, err)
		return nil, err
	}

	momentId, err := ret.LastInsertId()
	if err != nil {
		l.Logger.Errorf("LastInsertId error: %v", err)
		return nil, err
	}

	// 更新缓存
	threading.GoSafe(func() {
		l.addCache(context.Background(), in.UserId, now.Unix(), momentId)
	})

	return &pb.PublishResponse{MomentId: momentId}, nil
}

func (l *PublishLogic) addCache(ctx context.Context, userId, publishTime, momentId int64) {
	var score int64
	// 点赞数排序
	key := momentsKey(userId, types.SortLikeCount)
	b, err := l.svcCtx.BizRedis.ExistsCtx(ctx, key)
	if err != nil {
		logx.Errorf("ExistsCtx req: %v error: %v", key, err)
	}
	if b {
		score = 0
		logx.Errorf("ZaddCtx req: %v score: %v", key, score)
		_, err = l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.Itoa(int(momentId)))
		if err != nil {
			logx.Errorf("ZaddCtx req: %v error: %v", key, err)
		}
	}

	// 发布时间排序
	key = momentsKey(userId, types.SortPublishTime)
	b, err = l.svcCtx.BizRedis.ExistsCtx(ctx, key)
	if err != nil {
		logx.Errorf("ExistsCtx req: %v error: %v", key, err)
	}
	if b {
		score = publishTime
		logx.Errorf("ZaddCtx req: %v score: %v", key, score)
		_, err = l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.Itoa(int(momentId)))
		if err != nil {
			logx.Errorf("ZaddCtx req: %v error: %v", key, err)
		}
	}
}
