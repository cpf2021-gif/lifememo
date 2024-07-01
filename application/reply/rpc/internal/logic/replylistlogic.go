package logic

import (
	"context"
	"errors"
	"lifememo/application/reply/rpc/internal/code"
	"lifememo/application/reply/rpc/internal/model"
	"lifememo/application/reply/rpc/internal/types"

	"lifememo/application/reply/rpc/internal/svc"
	"lifememo/application/reply/rpc/pb1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyListLogic {
	return &ReplyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReplyListLogic) ReplyList(in *pb1.ReplyListRequest) (*pb1.ReplyListResponse, error) {
	resultParentReplyList := make([]*pb1.RootReplyItem, 0)

	var sortType string
	if in.SortType != types.SortPublishTime && in.SortType != types.SortLikeCount {
		return &pb1.ReplyListResponse{}, errors.New("sort type error")
	}

	if in.SortType == types.SortPublishTime {
		sortType = "create_time"
	} else {
		sortType = "like_num"
	}

	// 1. 找到对应target_id的所有根评论
	parentReplyList, err := l.svcCtx.ReplyModel.FindRootReplyByTargetId(in.TargetId, sortType)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, code.ReplyNotFound
		}
		l.Logger.Errorf("FindRootReplyByTargetId error: %v", err)
		return nil, err
	}

	for _, parentReply := range parentReplyList {
		resultSonReplyList := make([]*pb1.ReplyItem, 0)

		rootReply := &pb1.RootReplyItem{
			ReplyId:       parentReply.Id,
			BizId:         parentReply.BizId,
			TargetId:      parentReply.TargetId,
			ReplyUserId:   parentReply.ReplyUserId,
			BeReplyUserId: parentReply.BeReplyUserId,
			ParentId:      parentReply.ParentId,
			Content:       parentReply.Content,
			Status:        parentReply.Status,
			CreateTime:    parentReply.CreateTime.Unix(),
			UpdateTime:    parentReply.UpdateTime.Unix(),
			LikeNum:       parentReply.LikeNum,
		}

		// 2. 找到对应根评论的前2条评论
		replyList, err := l.svcCtx.ReplyModel.FindLimitReplyByParentId(parentReply.Id, 2)
		if err != nil && !errors.Is(err, model.ErrNotFound) {
			l.Logger.Errorf("FindReplyByParentId error: %v", err)
			return nil, err
		}

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
				LikeNum:       reply.LikeNum,
			}
			resultSonReplyList = append(resultSonReplyList, sonReply)
		}

		rootReply.ReplyItems = resultSonReplyList
		resultParentReplyList = append(resultParentReplyList, rootReply)
	}

	return &pb1.ReplyListResponse{
		RootReplyItems: resultParentReplyList,
	}, nil
}
