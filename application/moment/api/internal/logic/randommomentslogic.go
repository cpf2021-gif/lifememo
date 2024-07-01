package logic

import (
	"context"
	"lifememo/application/moment/rpc/pb"
	"lifememo/application/user/rpc/user"

	"lifememo/application/moment/api/internal/svc"
	"lifememo/application/moment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RandomMomentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRandomMomentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RandomMomentsLogic {
	return &RandomMomentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RandomMomentsLogic) RandomMoments() (*types.RandomMomentsResponse, error) {
	var ret []*types.MomentItem
	resp, err := l.svcCtx.MomentRPC.RandomMoments(l.ctx, &pb.RandomMomentsRequest{})
	if err != nil {
		return nil, err
	}

	uidSet := make(map[int64]struct{})
	uidNameMap := make(map[int64]string)

	for _, v := range resp.Moments {
		ret = append(ret, &types.MomentItem{
			Id:           v.Id,
			AuthorId:     v.AuthorId,
			Content:      v.Content,
			CommentCount: v.CommentCount,
			LikeCount:    v.LikeCount,
			PublishTime:  v.PublishTime,
		})

		uidSet[v.AuthorId] = struct{}{}
	}

	// 获取用户信息
	for k, _ := range uidSet {
		User, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{
			UserId: k,
		})
		if err != nil {
			return nil, err
		}

		uidNameMap[k] = User.Username
	}

	for _, v := range ret {
		v.AuthorName = uidNameMap[v.AuthorId]
	}

	return &types.RandomMomentsResponse{
		Moments: ret,
	}, nil
}
