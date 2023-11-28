package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MomentModel = (*customMomentModel)(nil)

type (
	// MomentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMomentModel.
	MomentModel interface {
		momentModel
		MomentsByUserId(ctx context.Context, userId, likeNum int64, status int, pubTime, sortField string, limit int) ([]*Moment, error)
		UpdateMomentStatus(ctx context.Context, id int64, status int) error
	}

	customMomentModel struct {
		*defaultMomentModel
	}
)

// NewMomentModel returns a model for the database table.
func NewMomentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MomentModel {
	return &customMomentModel{
		defaultMomentModel: newMomentModel(conn, c, opts...),
	}
}

func (m *customMomentModel) MomentsByUserId(ctx context.Context, userId, likeNum int64, status int, pubTime, sortField string, limit int) ([]*Moment, error) {
	var (
		err      error
		query    string
		anyField any
		articles []*Moment
	)
	if sortField == "like_num" {
		anyField = likeNum
		query = fmt.Sprintf("select "+momentRows+" from "+m.table+" where author_id=? and status=? and like_num < ? order by %s desc limit ?", sortField)
	} else {
		anyField = pubTime
		query = fmt.Sprintf("select "+momentRows+" from "+m.table+" where author_id=? and status=? and publish_time < ? order by %s desc limit ?", sortField)
	}
	err = m.QueryRowsNoCacheCtx(ctx, &articles, query, userId, status, anyField, limit)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (m *customMomentModel) UpdateMomentStatus(ctx context.Context, id int64, status int) error {
	lifeMemoMomentIdKey := fmt.Sprintf("%s%v", cacheLifememoMomentMomentIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("update " + m.table + " set status=? where id=?")
		return conn.ExecCtx(ctx, query, status, id)
	}, lifeMemoMomentIdKey)

	return err
}
