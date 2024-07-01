package logic

import (
	"context"
	"encoding/json"
	"lifememo/application/like/api/internal/code"
	"lifememo/application/like/rpc/like"
	"lifememo/application/moment/rpc/moment"
	"lifememo/application/reply/rpc/reply"
	"lifememo/pkg/xcode"

	"lifememo/application/like/api/internal/svc"
	"lifememo/application/like/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ThumbUpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThumbUpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbUpLogic {
	return &ThumbUpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThumbUpLogic) ThumbUp(req *types.ThumbUpRequest) (*types.ThumbUpResponse, error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("l.ctx.Value.error: %v", err)
		return nil, xcode.NoLogin
	}

	if req.BizId != types.MomentBizId && req.BizId != types.ReplyBizId {
		return nil, code.BizIdError
	}

	if req.ObjId <= 0 {
		return nil, code.ObjIdInvalid
	}

	if req.BizId == types.MomentBizId {
		_, err = l.svcCtx.MomentRPC.MomentDetail(l.ctx, &moment.MomentDetailRequest{
			MomentId: req.ObjId,
		})
		if err != nil {
			return nil, err
		}
	} else {
		_, err = l.svcCtx.ReplyRPC.FindReply(l.ctx, &reply.FindReplyRequest{
			ReplyId: req.ObjId,
		})
		if err != nil {
			return nil, err
		}
	}

	if req.LikeType != int32(types.Like) && req.LikeType != int32(types.UnLike) {
		return nil, code.TypeInvalid
	}

	resp, err := l.svcCtx.LikeRPC.IsThumbuped(l.ctx, &like.IsThumbupedRequest{
		BizId:  req.BizId,
		ObjId:  req.ObjId,
		UserId: userId,
	})

	if err != nil {
		l.Logger.Errorf("l.svcCtx.LikeRPC.IsThumbuped error: %v", err)
		return nil, err
	}

	if !resp.IsThumbuped && req.LikeType == int32(types.UnLike) {
		return nil, code.NotThumbuped
	} else if resp.IsThumbuped && req.LikeType == int32(types.Like) {
		return nil, code.AlreadyThumbuped
	} else {
		_, err := l.svcCtx.LikeRPC.Thumbup(l.ctx, &like.ThumbupRequest{
			BizId:    req.BizId,
			ObjId:    req.ObjId,
			UserId:   userId,
			LikeType: req.LikeType,
		})
		if err != nil {
			l.Logger.Errorf("l.svcCtx.LikeRPC.Thumbup error: %v", err)
			return nil, err
		}
		return &types.ThumbUpResponse{}, nil
	}
}
