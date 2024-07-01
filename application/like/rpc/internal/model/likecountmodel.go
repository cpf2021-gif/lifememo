package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LikeCountModel = (*customLikeCountModel)(nil)

type (
	// LikeCountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLikeCountModel.
	LikeCountModel interface {
		likeCountModel
		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		InsertWithSession(ctx context.Context, session sqlx.Session, data *LikeCount) (sql.Result, error)
		UpdateWithSession(ctx context.Context, session sqlx.Session, data *LikeCount) error
	}

	customLikeCountModel struct {
		*defaultLikeCountModel
	}
)

// NewLikeCountModel returns a model for the database table.
func NewLikeCountModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LikeCountModel {
	return &customLikeCountModel{
		defaultLikeCountModel: newLikeCountModel(conn, c, opts...),
	}
}

func (m *customLikeCountModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customLikeCountModel) InsertWithSession(ctx context.Context, session sqlx.Session, data *LikeCount) (sql.Result, error) {
	lifememoLikeLikeCountBizIdObjIdKey := fmt.Sprintf("%s%v:%v", cacheLifememoLikeLikeCountBizIdObjIdPrefix, data.BizId, data.ObjId)
	lifememoLikeLikeCountIdKey := fmt.Sprintf("%s%v", cacheLifememoLikeLikeCountIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, likeCountRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.BizId, data.ObjId, data.LikeNum)
	}, lifememoLikeLikeCountBizIdObjIdKey, lifememoLikeLikeCountIdKey)
	return ret, err
}

func (m *customLikeCountModel) UpdateWithSession(ctx context.Context, session sqlx.Session, newData *LikeCount) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	lifememoLikeLikeCountBizIdObjIdKey := fmt.Sprintf("%s%v:%v", cacheLifememoLikeLikeCountBizIdObjIdPrefix, data.BizId, data.ObjId)
	lifememoLikeLikeCountIdKey := fmt.Sprintf("%s%v", cacheLifememoLikeLikeCountIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, likeCountRowsWithPlaceHolder)
		return session.ExecCtx(ctx, query, newData.BizId, newData.ObjId, newData.LikeNum, newData.Id)
	}, lifememoLikeLikeCountBizIdObjIdKey, lifememoLikeLikeCountIdKey)
	return err
}
