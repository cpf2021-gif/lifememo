package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"lifememo/application/moment/mq/internal/svc"
	"lifememo/application/moment/mq/internal/types"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

type MomentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMomentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MomentLogic {
	return &MomentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MomentLogic) Consume(_, val string) error {
	var msg *types.CanalMsg
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		logx.Errorf("Consume value: %s, err: %v", val, err)
		return err
	}

	fmt.Println("Consume value: ", val)

	// add data to es
	if msg.Type == "INSERT" {

	}

	return l.AddEsData(msg)
}

func (l *MomentLogic) AddEsData(msg *types.CanalMsg) error {
	// add data to es
	if len(msg.Data) == 0 {
		return nil
	}

	var esData []*types.MomentEsMsg
	for _, d := range msg.Data {
		momentId, _ := strconv.ParseInt(d.ID, 10, 64)
		authorId, _ := strconv.ParseInt(d.AuthorID, 10, 64)
		status, _ := strconv.Atoi(d.Status)
		commentNum, _ := strconv.ParseInt(d.CommentNum, 10, 64)
		likeNum, _ := strconv.ParseInt(d.LikeNum, 10, 64)
		collectNum, _ := strconv.ParseInt(d.CollectNum, 10, 64)
		viewNum, _ := strconv.ParseInt(d.ViewNum, 10, 64)
		shareNum, _ := strconv.ParseInt(d.ShareNum, 10, 64)

		esData = append(esData, &types.MomentEsMsg{
			MomentID:    momentId,
			Content:     d.Content,
			AuthorID:    authorId,
			Status:      status,
			CommentNum:  commentNum,
			LikeNum:     likeNum,
			CollectNum:  collectNum,
			ViewNum:     viewNum,
			ShareNum:    shareNum,
			TagIds:      d.TagIds,
			PublishTime: d.PublishTime,
			CreateTime:  d.CreateTime,
			UpdateTime:  d.UpdateTime,
		})
	}

	return l.BatchUpsertToEs(l.ctx, esData)
}

func (l *MomentLogic) BatchUpsertToEs(ctx context.Context, esData []*types.MomentEsMsg) error {
	if len(esData) == 0 {
		return nil
	}

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: l.svcCtx.Es.Client,
		Index:  "moment-index",
	})
	if err != nil {
		return err
	}

	for _, d := range esData {
		v, err := json.Marshal(d)
		if err != nil {
			return err
		}

		err = bi.Add(ctx, esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: fmt.Sprintf("%d", d.MomentID),
			Body:       bytes.NewReader(v),
			OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, item2 esutil.BulkIndexerResponseItem) {
			},
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, item2 esutil.BulkIndexerResponseItem, err error) {
			},
		})

		if err != nil {
			return err
		}
	}

	return bi.Close(ctx)
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, NewMomentLogic(ctx, svcCtx)),
	}
}
