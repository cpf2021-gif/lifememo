package logic

import (
	"context"
	"errors"
	"lifememo/application/reply/rpc/internal/code"
	"lifememo/application/reply/rpc/internal/model"

	"lifememo/application/reply/rpc/internal/svc"
	"lifememo/application/reply/rpc/pb1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyDetailLogic {
	return &ReplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReplyDetailLogic) ReplyDetail(in *pb1.ReplyDetailRequest) (*pb1.ReplyDetailResponse, error) {
	// 1. 找到对应的根评论
	rootReply, err := l.svcCtx.ReplyModel.FindOne(l.ctx, in.ReplyId)
	if err != nil {
		l.Logger.Errorf("FindOne error: %v", err)
		if errors.Is(err, model.ErrNotFound) {
			return nil, code.ReplyNotFound
		}
		return nil, err
	}

	if rootReply.ParentId != 0 {
		return nil, code.ReplyIsNotRoot
	}

	// 2. 找到对应根评论的所有评论
	replyList, err := l.svcCtx.ReplyModel.FindLimitReplyByParentId(rootReply.Id, 10)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		l.Logger.Errorf("FindLimitReplyByParentId error: %v", err)
		return nil, err
	}

	// 3. 组装返回数据
	resultSonReplyList := make([]*pb1.ReplyItem, 0)
	for _, reply := range replyList {
		sonReply := &pb1.ReplyItem{
			ReplyId:       reply.Id,
			BizId:         reply.BizId,
			TargetId:      reply.TargetId,
			ReplyUserId:   reply.ReplyUserId,
			BeReplyUserId: reply.BeReplyUserId,
			ParentId:      reply.ParentId,
			Content:       reply.Content,
			Status:        reply.Status,
			CreateTime:    reply.CreateTime.Unix(),
			UpdateTime:    reply.UpdateTime.Unix(),
		}
		resultSonReplyList = append(resultSonReplyList, sonReply)
	}

	resultRootReply := &pb1.RootReplyItem{
		ReplyId:       rootReply.Id,
		BizId:         rootReply.BizId,
		TargetId:      rootReply.TargetId,
		ReplyUserId:   rootReply.ReplyUserId,
		BeReplyUserId: rootReply.BeReplyUserId,
		ParentId:      rootReply.ParentId,
		Content:       rootReply.Content,
		Status:        rootReply.Status,
		CreateTime:    rootReply.CreateTime.Unix(),
		UpdateTime:    rootReply.UpdateTime.Unix(),
		ReplyItems:    resultSonReplyList,
	}

	return &pb1.ReplyDetailResponse{
		RootReplyItem: resultRootReply,
	}, nil
}
