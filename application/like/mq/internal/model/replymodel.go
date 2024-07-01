package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ReplyModel = (*customReplyModel)(nil)

type (
	// ReplyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customReplyModel.
	ReplyModel interface {
		replyModel
		FindRootReplyByTargetId(targetId int64, sortType string) ([]*Reply, error)
		FindLimitReplyByParentId(parentId int64, limit int64) ([]*Reply, error)
	}

	customReplyModel struct {
		*defaultReplyModel
	}
)

// NewReplyModel returns a model for the database table.
func NewReplyModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ReplyModel {
	return &customReplyModel{
		defaultReplyModel: newReplyModel(conn, c, opts...),
	}
}

// FindRootReplyByTargetId 根据target_id找到所有根评论
func (m *customReplyModel) FindRootReplyByTargetId(targetId int64, sortType string) ([]*Reply, error) {
	var (
		err     error
		query   string
		replies []*Reply
	)

	if sortType == "create_time" {
		query = "select " + replyRows + " from " + m.table + " where target_id=? and parent_id=0 order by create_time desc"
	} else {
		query = "select " + replyRows + " from " + m.table + " where target_id=? and parent_id=0 order by like_num desc"
	}

	err = m.QueryRowsNoCache(&replies, query, targetId)
	if err != nil {
		return nil, err
	}

	return replies, nil
}

// FindLimitReplyByParentId 根据parent_id找到对应数量的评论评论
func (m *customReplyModel) FindLimitReplyByParentId(parentId int64, limit int64) ([]*Reply, error) {
	var (
		err     error
		query   string
		replies []*Reply
	)

	query = "select " + replyRows + " from " + m.table + " where parent_id=? order by create_time limit ?"
	err = m.QueryRowsNoCache(&replies, query, parentId, limit)
	if err != nil {
		return nil, err
	}

	return replies, nil
}
