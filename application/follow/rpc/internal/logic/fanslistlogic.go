package logic

import (
	"cmp"
	"context"
	"errors"
	"lifememo/application/follow/rpc/internal/model"
	"lifememo/application/follow/rpc/internal/svc"
	"lifememo/application/follow/rpc/internal/types"
	"lifememo/application/follow/rpc/pb"
	"slices"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
)

type FansListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFansListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FansListLogic {
	return &FansListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FansList 粉丝列表
func (l *FansListLogic) FansList(in *pb.FansListRequest) (*pb.FansListResponse, error) {
	if in.UserId <= 0 {
		return &pb.FansListResponse{}, errors.New("用户id非法")
	}
	if in.Id <= 0 {
		in.Id = 0
	}
	if in.PageSize <= 0 {
		in.PageSize = types.DefaultLimit
	}
	if in.Cursor <= 0 {
		in.Cursor = time.Now().Unix()
	}

	var (
		err            error
		isCache, isEnd bool
		lastId, cursor int64
		fansedUserIds  []int64
		sortCreateTime = time.Unix(in.Cursor, 0).Format("2006-01-02 15:04:05")
		fans           []*model.Follow
		curPage        []*pb.FansItem
	)

	fansUserIds, _ := l.cacheFanUserIds(l.ctx, in.UserId, in.Cursor, in.PageSize)
	if len(fansUserIds) > 0 {
		isCache = true
		if fansUserIds[len(fansUserIds)-1] == -1 {
			fansUserIds = fansUserIds[:len(fansUserIds)-1]
			isEnd = true
		}
		if len(fansUserIds) == 0 {
			return &pb.FansListResponse{
				IsEnd: true,
			}, nil
		}
		// MapReduce获取的数据是无序的，需要重新排序
		fans, err = l.FanListByIds(l.ctx, in.UserId, fansUserIds)

		cmpFunc := func(a, b *model.Follow) int {
			return cmp.Compare(b.CreateTime.Unix(), a.CreateTime.Unix())
		}
		slices.SortFunc(fans, cmpFunc)
		if err != nil {
			l.Logger.Errorf("[FansList] l.FanListByIds err: %v", err)
			return &pb.FansListResponse{}, err
		}
		for _, follow := range fans {
			fansedUserIds = append(fansedUserIds, follow.UserId)
			curPage = append(curPage, &pb.FansItem{
				Id:         follow.Id,
				FansUserId: follow.UserId,
				CreateTime: follow.CreateTime.Unix(),
			})
		}
	} else {
		fans, err = l.svcCtx.FollowModel.FansListByUserId(l.ctx, in.UserId, types.CacheMaxFollowCount, sortCreateTime)
		if err != nil {
			l.Logger.Errorf("[FansList] FollowModel.FansListByUserId err: %v", err)
			return nil, err
		}
		if len(fans) == 0 {
			return &pb.FansListResponse{
				IsEnd: true,
			}, nil
		}
		var firstPageFans []*model.Follow
		if len(fans) > int(in.PageSize) {
			firstPageFans = fans[:in.PageSize]
		} else {
			firstPageFans = fans
			isEnd = true
		}
		for _, follow := range firstPageFans {
			fansedUserIds = append(fansedUserIds, follow.UserId)
			curPage = append(curPage, &pb.FansItem{
				Id:         follow.Id,
				FansUserId: follow.UserId,
				CreateTime: follow.CreateTime.Unix(),
			})
		}
	}

	// 去重
	if len(curPage) > 0 {
		pageLast := curPage[len(curPage)-1]
		lastId = pageLast.Id
		cursor = pageLast.CreateTime
		if cursor < 0 {
			cursor = 0
		}
		for k, follow := range curPage {
			if follow.CreateTime == in.Cursor && follow.Id == in.Id {
				curPage = curPage[k:]
				break
			}
		}
	}
	fc, err := l.FollowCountListByIds(l.ctx, fansedUserIds)
	if err != nil {
		l.Logger.Errorf("[FansList] l.FollowCountListByIds err: %v", err)
	}
	uidFansCount := make(map[int64]int64)
	for _, f := range fc {
		uidFansCount[f.UserId] = f.FansCount
	}
	for _, cur := range curPage {
		cur.FansCount = uidFansCount[cur.FansUserId]
	}
	ret := &pb.FansListResponse{
		IsEnd:  isEnd,
		Cursor: cursor,
		Id:     lastId,
		Items:  curPage,
	}

	// 缓存
	if !isCache {
		threading.GoSafe(func() {
			if len(fans) < types.CacheMaxFollowCount && len(fans) > 0 {
				fans = append(fans, &model.Follow{
					UserId: -1,
				})
			}
			err = l.addCacheFans(context.Background(), in.UserId, fans)
			if err != nil {
				l.Logger.Errorf("[FansList] addCacheFans err: %v", err)
			}
		})
	}

	return ret, nil
}

func (l *FansListLogic) FanListByIds(ctx context.Context, userId int64, fansUserId []int64) ([]*model.Follow, error) {
	Fans, err := mr.MapReduce[int64, *model.Follow, []*model.Follow](func(source chan<- int64) {
		for _, uid := range fansUserId {
			if uid <= 0 {
				continue
			}
			source <- uid
		}
	}, func(uid int64, writer mr.Writer[*model.Follow], cancel func(error)) {
		p, err := l.svcCtx.FollowModel.FindOneByUserIdFollowedUserId(ctx, uid, userId)
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(p)
	}, func(pipe <-chan *model.Follow, writer mr.Writer[[]*model.Follow], cancel func(error)) {
		var Follows []*model.Follow
		for Follow := range pipe {
			Follows = append(Follows, Follow)
		}
		writer.Write(Follows)
	})
	if err != nil {
		return nil, err
	}

	return Fans, nil
}

func (l *FansListLogic) FollowCountListByIds(ctx context.Context, followedUserId []int64) ([]*model.FollowCount, error) {
	FollowCounts, err := mr.MapReduce[int64, *model.FollowCount, []*model.FollowCount](func(source chan<- int64) {
		for _, uid := range followedUserId {
			if uid <= 0 {
				continue
			}
			source <- uid
		}
	}, func(uid int64, writer mr.Writer[*model.FollowCount], cancel func(error)) {
		p, err := l.svcCtx.FollowCountModel.FindOneByUserId(ctx, uid)
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(p)
	}, func(pipe <-chan *model.FollowCount, writer mr.Writer[[]*model.FollowCount], cancel func(error)) {
		var FollowCounts []*model.FollowCount
		for FollowCount := range pipe {
			FollowCounts = append(FollowCounts, FollowCount)
		}
		writer.Write(FollowCounts)
	})
	if err != nil {
		return nil, err
	}

	return FollowCounts, nil
}

func (l *FansListLogic) cacheFanUserIds(ctx context.Context, userId, cursor, pageSize int64) ([]int64, error) {
	key := userFansKey(userId)
	b, err := l.svcCtx.BizRedis.ExistsCtx(ctx, key)
	if err != nil {
		logx.Errorf("[cacheFanUserIds] BizRedis.ExistsCtx err: %v", err)
	}
	if b {
		// 热点数据，更新过期时间
		err = l.svcCtx.BizRedis.ExpireCtx(ctx, key, userFollowExpire)
		if err != nil {
			logx.Errorf("[cacheFanUserIds] BizRedis.ExpireCtx err: %v", err)
		}
	}
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, cursor, 0, int(pageSize))
	if err != nil {
		logx.Errorf("[cacheFanUserIds] BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx err: %v", err)
		return nil, err
	}
	var uids []int64
	for _, pair := range pairs {
		uid, err := strconv.ParseInt(pair.Key, 10, 64)
		if err != nil {
			logx.Errorf("[cacheFanUserIds] strconv.ParseInt err: %v", err)
			continue
		}
		uids = append(uids, uid)
	}

	return uids, nil
}

func (l *FansListLogic) addCacheFans(ctx context.Context, userId int64, fans []*model.Follow) error {
	if len(fans) == 0 {
		return nil
	}
	key := userFansKey(userId)
	for _, follow := range fans {
		var score int64
		if follow.UserId == -1 {
			score = 0
		} else {
			score = follow.CreateTime.Unix()
		}
		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.FormatInt(follow.UserId, 10))
		if err != nil {
			logx.Errorf("[addCacheFans] BizRedis.ZaddCtx err: %v", err)
			return err
		}
	}

	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, userFollowExpire)
}
