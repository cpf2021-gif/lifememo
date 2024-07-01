package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LikeRecordModel = (*customLikeRecordModel)(nil)

type (
	// LikeRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLikeRecordModel.
	LikeRecordModel interface {
		likeRecordModel
		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		InsertWithSession(ctx context.Context, session sqlx.Session, data *LikeRecord) (sql.Result, error)
		UpdateWithSession(ctx context.Context, session sqlx.Session, data *LikeRecord) error
	}

	customLikeRecordModel struct {
		*defaultLikeRecordModel
	}
)

// NewLikeRecordModel returns a model for the database table.
func NewLikeRecordModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LikeRecordModel {
	return &customLikeRecordModel{
		defaultLikeRecordModel: newLikeRecordModel(conn, c, opts...),
	}
}

func (m *customLikeRecordModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customLikeRecordModel) InsertWithSession(ctx context.Context, session sqlx.Session, data *LikeRecord) (sql.Result, error) {
	lifememoLikeLikeRecordBizIdObjIdUserIdKey := fmt.Sprintf("%s%v:%v:%v", cacheLifememoLikeLikeRecordBizIdObjIdUserIdPrefix, data.BizId, data.ObjId, data.UserId)
	lifememoLikeLikeRecordIdKey := fmt.Sprintf("%s%v", cacheLifememoLikeLikeRecordIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, likeRecordRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.BizId, data.ObjId, data.UserId, data.LikeType)
	}, lifememoLikeLikeRecordBizIdObjIdUserIdKey, lifememoLikeLikeRecordIdKey)
	return ret, err
}

func (m *customLikeRecordModel) UpdateWithSession(ctx context.Context, session sqlx.Session, newData *LikeRecord) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	lifememoLikeLikeRecordBizIdObjIdUserIdKey := fmt.Sprintf("%s%v:%v:%v", cacheLifememoLikeLikeRecordBizIdObjIdUserIdPrefix, data.BizId, data.ObjId, data.UserId)
	lifememoLikeLikeRecordIdKey := fmt.Sprintf("%s%v", cacheLifememoLikeLikeRecordIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, likeRecordRowsWithPlaceHolder)
		return session.ExecCtx(ctx, query, newData.BizId, newData.ObjId, newData.UserId, newData.LikeType, newData.Id)
	}, lifememoLikeLikeRecordBizIdObjIdUserIdKey, lifememoLikeLikeRecordIdKey)
	return err
}
