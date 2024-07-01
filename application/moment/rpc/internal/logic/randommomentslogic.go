package logic

import (
	"context"
	"errors"

	"lifememo/application/moment/rpc/internal/svc"
	"lifememo/application/moment/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RandomMomentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRandomMomentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RandomMomentsLogic {
	return &RandomMomentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RandomMomentsLogic) RandomMoments(in *pb.RandomMomentsRequest) (*pb.RandomMomentsResponse, error) {
	moments, err := l.svcCtx.MomentModel.RandomMoments(l.ctx, 10)
	if err != nil {
		return &pb.RandomMomentsResponse{}, errors.New("获取随机动态失败")
	}

	var ret []*pb.MomentItem
	for _, moment := range moments {
		ret = append(ret, &pb.MomentItem{
			Id:           moment.Id,
			AuthorId:     moment.AuthorId,
			Content:      moment.Content,
			PublishTime:  moment.PublishTime.Unix(),
			LikeCount:    moment.LikeNum,
			CommentCount: moment.CommentNum,
		})
	}

	return &pb.RandomMomentsResponse{
		Moments: ret,
	}, nil
}
