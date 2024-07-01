package logic

import (
	"context"
	"errors"
	"lifememo/application/follow/rpc/internal/code"
	"lifememo/application/follow/rpc/internal/model"
	"lifememo/application/follow/rpc/internal/svc"
	"lifememo/application/follow/rpc/internal/types"
	"lifememo/application/follow/rpc/pb"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UnFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnFollowLogic {
	return &UnFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UnFollow 取消关注
func (l *UnFollowLogic) UnFollow(in *pb.UnFollowRequest) (*pb.UnFollowResponse, error) {
	if in.UserId <= 0 || in.FollowedUserId <= 0 {
		return &pb.UnFollowResponse{}, code.FollowUserIdInvalid
	}

	if in.UserId == in.FollowedUserId {
		return &pb.UnFollowResponse{}, code.DontFollowYourself
	}

	follow, err := l.svcCtx.FollowModel.FindOneByUserIdFollowedUserId(l.ctx, in.UserId, in.FollowedUserId)

	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return &pb.UnFollowResponse{}, err
	}

	if errors.Is(err, model.ErrNotFound) {
		return &pb.UnFollowResponse{}, code.NoFollowed
	}

	if follow != nil && follow.FollowStatus == types.FollowStatusUnfollow {
		return &pb.UnFollowResponse{}, code.UnFollowed
	}

	// 事务
	err = l.svcCtx.FollowCountModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 更新关注状态
		follow.FollowStatus = types.FollowStatusUnfollow
		err = l.svcCtx.FollowModel.UpdateWithSession(ctx, session, follow)
		if err != nil {
			return err
		}

		// 更新关注数, 粉丝数
		// 关注数 -1
		fc, err := l.svcCtx.FollowCountModel.FindOneByUserId(l.ctx, in.UserId)
		if err != nil {
			return err
		}
		err = l.svcCtx.FollowCountModel.UpdateFollowCount(l.ctx, session, fc.Id, fc.UserId, -1)
		if err != nil {
			return err
		}

		// 粉丝数 -1
		fc, err = l.svcCtx.FollowCountModel.FindOneByUserId(l.ctx, in.FollowedUserId)
		if err != nil {
			return err
		}
		err = l.svcCtx.FollowCountModel.UpdateFansCount(l.ctx, session, fc.Id, fc.UserId, -1)
		if err != nil {
			return err
		}

		return err
	})

	if err != nil {
		l.Logger.Errorf("[UnFollow] UnFollow Transaction error: %v", err)
		return &pb.UnFollowResponse{}, err
	}

	// 删除缓存
	// followList
	followExist, err := l.svcCtx.BizRedis.ExistsCtx(l.ctx, userFollowKey(in.UserId))
	if err != nil {
		l.Logger.Errorf("[UnFollow] ExistsCtx error: %v", err)
		return &pb.UnFollowResponse{}, err
	}
	if followExist {
		_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, userFollowKey(in.UserId), strconv.FormatInt(in.FollowedUserId, 10))
		if err != nil {
			l.Logger.Errorf("[UnFollow] ZremCtx error: %v", err)
			return &pb.UnFollowResponse{}, err
		}
	}
	// fansList
	fansExist, err := l.svcCtx.BizRedis.ExistsCtx(l.ctx, userFansKey(in.FollowedUserId))
	if err != nil {
		l.Logger.Errorf("[UnFollow] ExistsCtx error: %v", err)
		return &pb.UnFollowResponse{}, err
	}
	if fansExist {
		_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, userFansKey(in.FollowedUserId), strconv.FormatInt(in.UserId, 10))
		if err != nil {
			l.Logger.Errorf("[UnFollow] ZremCtx error: %v", err)
			return &pb.UnFollowResponse{}, err
		}
	}

	return &pb.UnFollowResponse{}, nil
}
