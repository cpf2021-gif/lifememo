package logic

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"lifememo/application/follow/rpc/internal/model"
	"lifememo/application/follow/rpc/internal/types"
	"slices"
	"strconv"
	"time"

	"lifememo/application/follow/rpc/internal/svc"
	"lifememo/application/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
)

const userFollowExpire = 3600 * 24 * 2

type FollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FollowList 关注列表
func (l *FollowListLogic) FollowList(in *pb.FollowListRequest) (*pb.FollowListResponse, error) {
	if in.UserId <= 0 {
		return &pb.FollowListResponse{}, errors.New("用户id非法")
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
		err             error
		isCache, isEnd  bool
		lastId, cursor  int64
		followedUserIds []int64
		sortCreateTime  = time.Unix(in.Cursor, 0).Format("2006-01-02 15:04:05")
		follows         []*model.Follow
		curPage         []*pb.FollowItem
	)

	followUserIds, _ := l.cacheFollowUserIds(l.ctx, in.UserId, in.Cursor, in.PageSize)
	if len(followUserIds) > 0 {
		isCache = true
		if followUserIds[len(followUserIds)-1] == -1 {
			followUserIds = followUserIds[:len(followUserIds)-1]
			isEnd = true
		}
		if len(followUserIds) == 0 {
			return &pb.FollowListResponse{
				IsEnd: true,
			}, nil
		}
		// MapReduce获取的数据是无序的，需要重新排序
		follows, err = l.FollowListByIds(l.ctx, in.UserId, followUserIds)

		cmpFunc := func(a, b *model.Follow) int {
			return cmp.Compare(b.CreateTime.Unix(), a.CreateTime.Unix())
		}
		slices.SortFunc(follows, cmpFunc)

		if err != nil {
			l.Logger.Errorf("[FollowList] FollowListByIds err: %v", err)
			return &pb.FollowListResponse{}, err
		}
		for _, follow := range follows {
			followedUserIds = append(followedUserIds, follow.FollowedUserId)
			curPage = append(curPage, &pb.FollowItem{
				Id:             follow.Id,
				FollowedUserId: follow.FollowedUserId,
				CreateTime:     follow.CreateTime.Unix(),
			})
		}
	} else {
		follows, err = l.svcCtx.FollowModel.FollowListByUserId(l.ctx, in.UserId, types.CacheMaxFollowCount, sortCreateTime)
		if err != nil {
			l.Logger.Errorf("[FollowList] FollowModel.FollowListByUserId err: %v", err)
			return nil, err
		}
		if len(follows) == 0 {
			return &pb.FollowListResponse{
				IsEnd: true,
			}, nil
		}
		var firstPageFollows []*model.Follow
		if len(follows) > int(in.PageSize) {
			firstPageFollows = follows[:in.PageSize]
		} else {
			firstPageFollows = follows
			isEnd = true
		}
		for _, follow := range firstPageFollows {
			followedUserIds = append(followedUserIds, follow.FollowedUserId)
			curPage = append(curPage, &pb.FollowItem{
				Id:             follow.Id,
				FollowedUserId: follow.FollowedUserId,
				CreateTime:     follow.CreateTime.Unix(),
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
	fc, err := l.FollowCountListByIds(l.ctx, followedUserIds)
	if err != nil {
		l.Logger.Errorf("[FollowList] FollowCountListByIds err: %v", err)
	}
	uidFansCount := make(map[int64]int64)
	for _, f := range fc {
		uidFansCount[f.UserId] = f.FansCount
	}
	for _, cur := range curPage {
		cur.FansCount = uidFansCount[cur.FollowedUserId]
	}
	ret := &pb.FollowListResponse{
		IsEnd:  isEnd,
		Cursor: cursor,
		Id:     lastId,
		Items:  curPage,
	}

	// 缓存
	if !isCache {
		threading.GoSafe(func() {
			if len(follows) < types.CacheMaxFollowCount && len(follows) > 0 {
				follows = append(follows, &model.Follow{FollowedUserId: -1})
			}
			err = l.addCacheFollow(context.Background(), in.UserId, follows)
			if err != nil {
				l.Logger.Errorf("[FollowList] addCacheFollow err: %v", err)
			}
		})
	}

	return ret, nil
}

func (l *FollowListLogic) FollowListByIds(ctx context.Context, userId int64, followedUserId []int64) ([]*model.Follow, error) {
	Follows, err := mr.MapReduce[int64, *model.Follow, []*model.Follow](func(source chan<- int64) {
		for _, uid := range followedUserId {
			if uid <= 0 {
				continue
			}
			source <- uid
		}
	}, func(uid int64, writer mr.Writer[*model.Follow], cancel func(error)) {
		p, err := l.svcCtx.FollowModel.FindOneByUserIdFollowedUserId(ctx, userId, uid)
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

	return Follows, nil
}

func (l *FollowListLogic) FollowCountListByIds(ctx context.Context, followedUserId []int64) ([]*model.FollowCount, error) {
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

func (l *FollowListLogic) cacheFollowUserIds(ctx context.Context, userId, cursor, pageSize int64) ([]int64, error) {
	key := userFollowKey(userId)
	b, err := l.svcCtx.BizRedis.ExistsCtx(ctx, key)
	if err != nil {
		logx.Errorf("[cacheFollowUserIds] BizRedis.ExistsCtx err: %v", err)
	}
	if b {
		// 热点数据，更新过期时间
		err = l.svcCtx.BizRedis.ExpireCtx(ctx, key, userFollowExpire)
		if err != nil {
			logx.Errorf("[cacheFollowUserIds] BizRedis.ExpireCtx err: %v", err)
		}
	}
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, cursor, 0, int(pageSize))
	if err != nil {
		logx.Errorf("[cacheFollowUserIds] BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx err: %v", err)
		return nil, err
	}
	var uids []int64
	for _, pair := range pairs {
		uid, err := strconv.ParseInt(pair.Key, 10, 64)
		if err != nil {
			logx.Errorf("[cacheFollowUserIds] strconv.ParseInt err: %v", err)
			continue
		}
		uids = append(uids, uid)
	}

	return uids, nil
}

func (l *FollowListLogic) addCacheFollow(ctx context.Context, userId int64, follows []*model.Follow) error {
	if len(follows) == 0 {
		return nil
	}
	key := userFollowKey(userId)
	for _, follow := range follows {
		var score int64
		if follow.FollowedUserId == -1 {
			score = 0
		} else {
			score = follow.CreateTime.Unix()
		}
		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.FormatInt(follow.FollowedUserId, 10))
		if err != nil {
			logx.Errorf("[addCacheFollow] BizRedis.ZaddCtx err: %v", err)
			return err
		}
	}

	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, userFollowExpire)
}

func userFollowKey(userId int64) string {
	return fmt.Sprintf("biz#user#follow#%d", userId)
}

func userFansKey(userId int64) string {
	return fmt.Sprintf("biz#user#fans#%d", userId)
}
