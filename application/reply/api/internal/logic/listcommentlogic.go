package logic

import (
	"context"
	"lifememo/application/reply/api/internal/code"
	"lifememo/application/reply/rpc/reply"

	"lifememo/application/reply/api/internal/svc"
	"lifememo/application/reply/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentLogic {
	return &ListCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCommentLogic) ListComment(req *types.CommentListRequest) (*types.CommentListResponse, error) {
	if req.BizId != types.MomentBizId {
		return nil, code.ErrInvalidBizId
	}

	if req.TargetId <= 0 {
		return nil, code.ErrInvalidTargetId
	}

	resp, err := l.svcCtx.ReplyRPC.ReplyList(l.ctx, &reply.ReplyListRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
		SortType: req.SortType,
		ReplyId:  req.ReplyId,
	})

	if err != nil {
		return nil, err
	}

	var rootCommentList []*types.RootReplyItem
	for _, item := range resp.RootReplyItems {
		var ReplyItems []*types.ReplyItem
		for _, replyItem := range item.ReplyItems {
			ReplyItems = append(ReplyItems, &types.ReplyItem{
				ReplyId:         replyItem.ReplyId,
				BizId:           replyItem.BizId,
				TargetId:        replyItem.TargetId,
				ReplyUserId:     replyItem.ReplyUserId,
				BeRepliedUserId: replyItem.BeReplyUserId,
				Content:         replyItem.Content,
				Status:          replyItem.Status,
				CreateTime:      replyItem.CreateTime,
				UpdateTime:      replyItem.UpdateTime,
				LikeNum:         replyItem.LikeNum,
			})
		}
		rootCommentList = append(rootCommentList, &types.RootReplyItem{
			ReplyId:         item.ReplyId,
			BizId:           item.BizId,
			TargetId:        item.TargetId,
			ReplyUserId:     item.ReplyUserId,
			BeRepliedUserId: item.BeReplyUserId,
			Content:         item.Content,
			Status:          item.Status,
			CreateTime:      item.CreateTime,
			UpdateTime:      item.UpdateTime,
			LikeNum:         item.LikeNum,
			ReplyItems:      ReplyItems,
		})
	}

	return &types.CommentListResponse{
		ReplyItems: rootCommentList,
	}, nil
}
