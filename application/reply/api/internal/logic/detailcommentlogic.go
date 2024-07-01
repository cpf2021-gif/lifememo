package logic

import (
	"context"
	"lifememo/application/reply/rpc/reply"

	"lifememo/application/reply/api/internal/svc"
	"lifememo/application/reply/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailCommentLogic {
	return &DetailCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailCommentLogic) DetailComment(req *types.CommentDetailRequest) (*types.CommentDetailResponse, error) {
	_, err := l.svcCtx.ReplyRPC.FindReply(l.ctx, &reply.FindReplyRequest{
		ReplyId: req.ReplyId,
	})

	if err != nil {
		return nil, err
	}

	resp, err := l.svcCtx.ReplyRPC.ReplyDetail(l.ctx, &reply.ReplyDetailRequest{
		ReplyId: req.ReplyId,
	})

	if err != nil {
		return nil, err
	}

	var ReplyItem []*types.ReplyItem
	for _, item := range resp.RootReplyItem.ReplyItems {
		ReplyItem = append(ReplyItem, &types.ReplyItem{
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
		})
	}

	var rootReplyItem *types.RootReplyItem
	rootReplyItem = &types.RootReplyItem{
		ReplyId:         resp.RootReplyItem.ReplyId,
		BizId:           resp.RootReplyItem.BizId,
		TargetId:        resp.RootReplyItem.TargetId,
		ReplyUserId:     resp.RootReplyItem.ReplyUserId,
		BeRepliedUserId: resp.RootReplyItem.BeReplyUserId,
		Content:         resp.RootReplyItem.Content,
		Status:          resp.RootReplyItem.Status,
		CreateTime:      resp.RootReplyItem.CreateTime,
		UpdateTime:      resp.RootReplyItem.UpdateTime,
		LikeNum:         resp.RootReplyItem.LikeNum,
		ReplyItems:      ReplyItem,
	}

	return &types.CommentDetailResponse{
		ReplyItem: rootReplyItem,
	}, nil
}
