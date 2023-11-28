package logic

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strconv"
	"time"

	"lifememo/application/moment/rpc/internal/code"
	"lifememo/application/moment/rpc/internal/model"
	"lifememo/application/moment/rpc/internal/svc"
	"lifememo/application/moment/rpc/internal/types"
	"lifememo/application/moment/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
)

const (
	prefixMoments = "biz#moments#%d#%d"
	momentsExpire = 3600 * 24 * 2
)

type MomentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMomentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MomentsLogic {
	return &MomentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MomentsLogic) Moments(in *pb.MomentsRequest) (*pb.MomentsResponse, error) {
	if in.SortType != types.SortLikeCount && in.SortType != types.SortPublishTime {
		return nil, code.SortTypeInvalid
	}
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor == 0 {
		if in.SortType == types.SortPublishTime {
			in.Cursor = time.Now().Unix()
		} else {
			in.Cursor = types.DefaultSortLikeCursor
		}
	}

	var (
		sortField       string
		sortLikeNum     int64
		sortPublishTime string
	)
	if in.SortType == types.SortLikeCount {
		sortField = "like_num"
		sortLikeNum = in.Cursor
	} else {
		sortField = "publish_time"
		sortPublishTime = time.Unix(in.Cursor, 0).Format("2006-01-02 15:04:05")
	}

	var (
		err            error
		isCache, isEnd bool
		lastId, cursor int64
		curPage        []*pb.MomentItem
		moments        []*model.Moment
	)

	momentIds, _ := l.CacheMoments(l.ctx, in.UserId, in.Cursor, in.PageSize, in.SortType)
	if len(momentIds) > 0 {
		isCache = true
		if momentIds[len(momentIds)-1] == -1 {
			isEnd = true
		}

		moments, err = l.MomentsByIds(l.ctx, momentIds)
		if err != nil {
			return nil, err
		}

		// mapreduce获取到的数据顺序可能会乱，需要重新排序
		var cmpFunc func(a, b *model.Moment) int
		if sortField == "like_num" {
			cmpFunc = func(a, b *model.Moment) int {
				return cmp.Compare(b.LikeNum, a.LikeNum)
			}
		} else {
			cmpFunc = func(a, b *model.Moment) int {
				return cmp.Compare(b.PublishTime.Unix(), a.PublishTime.Unix())
			}
		}
		slices.SortFunc(moments, cmpFunc)

		for _, moment := range moments {
			curPage = append(curPage, &pb.MomentItem{
				Id:           moment.Id,
				Content:      moment.Content,
				CommentCount: moment.CommentNum,
				LikeCount:    moment.LikeNum,
				PublishTime:  moment.PublishTime.Unix(),
				AuthorId:     moment.AuthorId,
			})
		}

	} else {
		v, err, _ := l.svcCtx.SingleFlightGroup.Do(fmt.Sprintf("MomentsByUserId:%d:%d", in.UserId, in.SortType), func() (interface{}, error) {
			return l.svcCtx.MomentModel.MomentsByUserId(l.ctx, in.UserId, sortLikeNum, types.MomentStatusVisible, sortPublishTime, sortField, int(in.PageSize))
		})
		if err != nil {
			logx.Errorf("MomentsByUserId userId: %d sortType: %d error: %v", in.UserId, in.SortType, err)
			return nil, err
		}
		if v == nil {
			return &pb.MomentsResponse{
				IsEnd: true,
			}, nil
		}
		moments = v.([]*model.Moment)
		var firstPageMoments []*model.Moment
		if len(moments) > int(in.PageSize) {
			firstPageMoments = moments[:int(in.PageSize)]
		} else {
			firstPageMoments = moments
			isEnd = true
		}
		for _, moment := range firstPageMoments {
			curPage = append(curPage, &pb.MomentItem{
				Id:           moment.Id,
				Content:      moment.Content,
				CommentCount: moment.CommentNum,
				LikeCount:    moment.LikeNum,
				PublishTime:  moment.PublishTime.Unix(),
				AuthorId:     moment.AuthorId,
			})
		}
	}

	if len(curPage) > 0 {
		pageLast := curPage[len(curPage)-1]
		lastId = pageLast.Id
		if in.SortType == types.SortLikeCount {
			cursor = pageLast.LikeCount
		} else {
			cursor = pageLast.PublishTime
		}
		if cursor < 0 {
			cursor = 0
		}
		for k, moment := range curPage {
			if in.SortType == types.SortLikeCount {
				if moment.LikeCount == in.Cursor && moment.Id == in.MomentId {
					curPage = curPage[k:]
					break
				}
			} else {
				if moment.PublishTime == in.Cursor && moment.Id == in.MomentId {
					curPage = curPage[k:]
					break
				}
			}
		}
	}

	ret := &pb.MomentsResponse{
		IsEnd:    isEnd,
		Cursor:   cursor,
		MomentId: lastId,
		Moments:  curPage,
	}

	if !isCache {
		threading.GoSafe(func() {
			if len(moments) < types.DefaultLimit && len(moments) > 0 {
				moments = append(moments, &model.Moment{Id: -1})
			}
			err = l.addCacheMoments(context.Background(), moments, in.UserId, in.SortType)
			if err != nil {
				logx.Errorf("addCacheMoments error: %v", err)
			}
		})
	}

	return ret, nil
}

func (l *MomentsLogic) MomentsByIds(ctx context.Context, momentIds []int64) ([]*model.Moment, error) {
	Moments, err := mr.MapReduce[int64, *model.Moment, []*model.Moment](func(source chan<- int64) {
		for _, mid := range momentIds {
			if mid == -1 {
				continue
			}
			source <- mid
		}
	}, func(id int64, writer mr.Writer[*model.Moment], cancel func(error)) {
		p, err := l.svcCtx.MomentModel.FindOne(ctx, id)
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(p)
	}, func(pipe <-chan *model.Moment, writer mr.Writer[[]*model.Moment], cancel func(error)) {
		var Moments []*model.Moment
		for Moment := range pipe {
			Moments = append(Moments, Moment)
		}
		writer.Write(Moments)
	})
	if err != nil {
		return nil, err
	}

	return Moments, nil
}

func (l *MomentsLogic) CacheMoments(ctx context.Context, uid, cursor, pageSize int64, sortType int32) ([]int64, error) {
	key := momentsKey(uid, sortType)
	b, err := l.svcCtx.BizRedis.ExistsCtx(ctx, key)
	if err != nil {
		logx.Errorf("ExistsCtx key: %s error: %v", key, err)
	}
	if b {
		// 对热数据进行续期
		err = l.svcCtx.BizRedis.ExpireCtx(ctx, key, momentsExpire)
		if err != nil {
			logx.Errorf("ExpireCtx key: %s error: %v", key, err)
		}
	}
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, cursor, 0, int(pageSize))
	if err != nil {
		logx.Errorf("ZrevrangebyscoreWithScoresAndLimitCtx key: %s error: %v", key, err)
		return nil, err
	}
	var ids []int64
	for _, pair := range pairs {
		id, err := strconv.ParseInt(pair.Key, 10, 64)
		if err != nil {
			logx.Errorf("strconv.ParseInt key: %s error: %v", pair.Key, err)
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (l *MomentsLogic) addCacheMoments(ctx context.Context, moments []*model.Moment, userId int64, sortType int32) error {
	if len(moments) == 0 {
		return nil
	}
	key := momentsKey(userId, sortType)
	for _, moment := range moments {
		var score int64
		if sortType == types.SortLikeCount {
			score = moment.LikeNum
		} else if sortType == types.SortPublishTime {
			score = moment.PublishTime.Unix()
		}
		if score < 0 {
			score = 0
		}
		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.Itoa(int(moment.Id)))
		if err != nil {
			logx.Errorf("addCacheMoments error: %v", err)
			return err
		}
	}

	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, momentsExpire)
}

func momentsKey(userId int64, sortType int32) string {
	return fmt.Sprintf(prefixMoments, userId, sortType)
}
