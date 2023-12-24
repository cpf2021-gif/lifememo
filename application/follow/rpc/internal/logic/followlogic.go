package logic

import (
	"context"
	"lifememo/application/follow/rpc/internal/model"
	"lifememo/application/follow/rpc/internal/types"
	"strconv"
	"time"

	"lifememo/application/follow/rpc/internal/svc"
	"lifememo/application/follow/rpc/pb"

	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Follow 关注
func (l *FollowLogic) Follow(in *pb.FollowRequest) (*pb.FollowResponse, error) {
	if in.UserId <= 0 || in.FollowedUserId <= 0 {
		return &pb.FollowResponse{}, errors.New("参数错误")
	}

	if in.UserId == in.FollowedUserId {
		return &pb.FollowResponse{}, errors.New("不能对自己进行操作")
	}

	follow, err := l.svcCtx.FollowModel.FindOneByUserIdFollowedUserId(l.ctx, in.UserId, in.FollowedUserId)

	if err != nil && err != model.ErrNotFound {
		return &pb.FollowResponse{}, err
	}

	if follow != nil && follow.FollowStatus == types.FollowStatusFollow {
		return &pb.FollowResponse{}, nil
	}

	// 事务
	err = l.svcCtx.FollowCountModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if err == model.ErrNotFound {
			_, err = l.svcCtx.FollowModel.InsertWithSession(l.ctx, session, &model.Follow{
				UserId:         in.UserId,
				FollowedUserId: in.FollowedUserId,
				FollowStatus:   types.FollowStatusFollow,
			})

			if err != nil {
				return err
			}
		} else {
			// 更新关注状态
			follow.FollowStatus = types.FollowStatusFollow
			err = l.svcCtx.FollowModel.UpdateWithSession(l.ctx, session, follow)
			if err != nil {
				return err
			}
		}

		// 更新关注数, 粉丝数
		// 关注数 +1
		fc, err := l.svcCtx.FollowCountModel.FindOneByUserId(l.ctx, in.UserId)
		if err != nil && err != model.ErrNotFound {
			return err
		} else if err == model.ErrNotFound {
			// 不存在则插入
			_, err = l.svcCtx.FollowCountModel.InsertWithSession(l.ctx, session, &model.FollowCount{
				UserId:      in.UserId,
				FollowCount: 1,
			})
			if err != nil {
				return err
			}
		} else {
			// 存在则更新
			err = l.svcCtx.FollowCountModel.UpdateFollowCount(l.ctx, session, fc.Id, fc.UserId, 1)
			if err != nil {
				return err
			}
		}

		// 粉丝数 +1
		fc, err = l.svcCtx.FollowCountModel.FindOneByUserId(l.ctx, in.FollowedUserId)
		if err != nil && !errors.Is(err, model.ErrNotFound) {
			return err
		} else if errors.Is(err, model.ErrNotFound) {
			// 不存在则插入
			_, err = l.svcCtx.FollowCountModel.InsertWithSession(l.ctx, session, &model.FollowCount{
				UserId:    in.FollowedUserId,
				FansCount: 1,
			})
			if err != nil {
				return err
			}
		} else {
			// 存在则更新
			err = l.svcCtx.FollowCountModel.UpdateFansCount(l.ctx, session, fc.Id, fc.UserId, 1)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		l.Logger.Errorf("[Follow] Follow Transaction error: %v", err)
		return &pb.FollowResponse{}, err
	}

	// 更新followList缓存
	followExist, err := l.svcCtx.BizRedis.ExistsCtx(l.ctx, userFollowKey(in.UserId))
	if err != nil {
		l.Logger.Errorf("[Follow] BizRedis.ExistsCtx error: %v", err)
		return &pb.FollowResponse{}, err
	}
	if followExist {
		// 添加缓存
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, userFollowKey(in.UserId), time.Now().Unix(), strconv.FormatInt(in.FollowedUserId, 10))
		if err != nil {
			l.Logger.Errorf("[Follow] BizRedis.ZaddCtx error: %v", err)
			return &pb.FollowResponse{}, err
		}
		// 删除多余的缓存
		_, err = l.svcCtx.BizRedis.ZremrangebyrankCtx(l.ctx, userFollowKey(in.UserId), 0, -(types.CacheMaxFollowCount + 1))
		if err != nil {
			l.Logger.Errorf("[Follow] BizRedis.ZremrangebyrankCtx error: %v", err)
		}
	}

	// 更新fansList缓存
	fansExist, err := l.svcCtx.BizRedis.ExistsCtx(l.ctx, userFansKey(in.FollowedUserId))
	if err != nil {
		l.Logger.Errorf("[Follow] BizRedis.ExistsCtx error: %v", err)
		return &pb.FollowResponse{}, err
	}
	if fansExist {
		// 添加缓存
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, userFansKey(in.FollowedUserId), time.Now().Unix(), strconv.FormatInt(in.UserId, 10))
		if err != nil {
			l.Logger.Errorf("[Follow] BizRedis.ZaddCtx error: %v", err)
			return &pb.FollowResponse{}, err
		}
		// 删除多余的缓存
		_, err = l.svcCtx.BizRedis.ZremrangebyrankCtx(l.ctx, userFansKey(in.FollowedUserId), 0, -(types.CacheMaxFansCount + 1))
		if err != nil {
			l.Logger.Errorf("[Follow] BizRedis.ZremrangebyrankCtx error: %v", err)
		}
	}

	return &pb.FollowResponse{}, nil
}
