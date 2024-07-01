package logic

import (
	"context"
	"errors"
	"lifememo/application/reply/rpc/internal/model"
	"lifememo/application/reply/rpc/internal/svc"
	"lifememo/application/reply/rpc/pb1"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishLogic) Publish(in *pb1.PublishRequest) (*pb1.PublishResponse, error) {
	reply := &model.Reply{}
	err := copier.Copy(reply, in)
	if err != nil {
		l.Logger.Error(err)
		return nil, err
	}

	result, err := l.svcCtx.ReplyModel.Insert(l.ctx, reply)
	if err != nil {
		l.Logger.Error(err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		l.Logger.Error(err)
		return nil, err
	}

	// 更新reply_count
	rc, err := l.svcCtx.ReplyCountModel.FindOneByBizIdTargetId(l.ctx, in.BizId, in.TargetId)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		l.Logger.Errorf("[PublishLogic] ReplyCountModel.FindOneByBizIdTargetId err: %v", err)
	} else if errors.Is(err, model.ErrNotFound) {
		// 不存在则创建
		rc = &model.ReplyCount{
			BizId:    in.BizId,
			TargetId: in.TargetId,
			ReplyNum: 1,
		}
		if in.ParentId == 0 {
			rc.ReplyRootNum = 1
		}
		_, err = l.svcCtx.ReplyCountModel.Insert(l.ctx, rc)
		if err != nil {
			l.Logger.Errorf("[PublishLogic] ReplyCountModel.Insert err: %v", err)
		}
	} else {
		// 存在则更新
		rc.ReplyNum++
		if in.ParentId == 0 {
			rc.ReplyRootNum++
		}
		err = l.svcCtx.ReplyCountModel.Update(l.ctx, rc)
		if err != nil {
			l.Logger.Errorf("[PublishLogic] ReplyCountModel.Update err: %v", err)
		}
	}

	// 异步更新动态回复数
	threading.GoSafe(func() {
		if in.BizId == "moment" {
			err = l.svcCtx.MomentCommentKqPusherClient.Push(strconv.FormatInt(in.TargetId, 10))
			if err != nil {
				l.Logger.Errorf("[PublishLogic] MomentCommentKqPusherClient.Push err: %v", err)
				return
			}
		}
	})

	return &pb1.PublishResponse{
		ReplyId: id,
	}, nil
}
