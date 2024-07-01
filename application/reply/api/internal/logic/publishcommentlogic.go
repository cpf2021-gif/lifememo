package logic

import (
	"context"
	"encoding/json"
	"lifememo/application/moment/rpc/moment"
	"lifememo/application/reply/api/internal/code"
	"lifememo/application/reply/rpc/reply"
	"lifememo/application/user/rpc/user"
	"lifememo/pkg/xcode"

	"lifememo/application/reply/api/internal/svc"
	"lifememo/application/reply/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishCommentLogic {
	return &PublishCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishCommentLogic) PublishComment(req *types.CommentPublishRequest) (*types.CommentPublishResponse, error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("l.ctx.Value.error: %v", err)
		return nil, xcode.NoLogin
	}

	if req.BizId != types.MomentBizId {
		return nil, code.ErrInvalidBizId
	}
	if req.TargetId <= 0 {
		return nil, code.ErrInvalidTargetId
	}
	if req.ParentId < 0 {
		return nil, code.ErrInvalidParentId
	}
	if req.BeRepliedId == userId || req.BeRepliedId < 0 {
		return nil, code.ErrInvalidBeRepliedId
	}

	if req.BeRepliedId > 0 {
		_, err = l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{
			UserId: req.BeRepliedId,
		})
		if err != nil {
			return nil, code.BeRepliedUserNotFound
		}
	}

	_, err = l.svcCtx.MomentRPC.MomentDetail(l.ctx, &moment.MomentDetailRequest{
		MomentId: req.TargetId,
	})

	if err != nil {
		return nil, err
	}

	resp, err := l.svcCtx.ReplyRPC.Publish(l.ctx, &reply.PublishRequest{
		BizId:         req.BizId,
		TargetId:      req.TargetId,
		ReplyUserId:   userId,
		BeReplyUserId: req.BeRepliedId,
		Content:       req.Content,
		ParentId:      req.ParentId,
		Status:        req.Status,
	})

	if err != nil {
		return nil, err
	}

	return &types.CommentPublishResponse{
		Id: resp.ReplyId,
	}, nil
}
