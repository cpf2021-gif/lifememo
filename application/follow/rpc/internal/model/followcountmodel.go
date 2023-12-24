package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FollowCountModel = (*customFollowCountModel)(nil)

type (
	// FollowCountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFollowCountModel.
	FollowCountModel interface {
		followCountModel
		UpdateFollowCount(ctx context.Context, session sqlx.Session, id, userId, opt int64) error
		UpdateFansCount(ctx context.Context, session sqlx.Session, id, userId, opt int64) error
		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		InsertWithSession(ctx context.Context, session sqlx.Session, data *FollowCount) (sql.Result, error)
	}

	customFollowCountModel struct {
		*defaultFollowCountModel
	}
)

// NewFollowCountModel returns a model for the database table.
func NewFollowCountModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FollowCountModel {
	return &customFollowCountModel{
		defaultFollowCountModel: newFollowCountModel(conn, c, opts...),
	}
}

// InsertWithSession insert a record into table with session
func (m *customFollowCountModel) InsertWithSession(ctx context.Context, session sqlx.Session, data *FollowCount) (sql.Result, error) {
	lifememoFollowFollowCountIdKey := fmt.Sprintf("%s%v", cacheLifememoFollowFollowCountIdPrefix, data.Id)
	lifememoFollowFollowCountUserIdKey := fmt.Sprintf("%s%v", cacheLifememoFollowFollowCountUserIdPrefix, data.UserId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, followCountRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.UserId, data.FollowCount, data.FansCount)
	}, lifememoFollowFollowCountIdKey, lifememoFollowFollowCountUserIdKey)
	return ret, err
}

// UpdateFollowCount 更新关注数 1 增加 -1 减少
func (m *customFollowCountModel) UpdateFollowCount(ctx context.Context, session sqlx.Session, id, userId, opt int64) error {
	lifememoFollowFollowCountIdKey := fmt.Sprintf("%s%v", cacheLifememoFollowFollowCountIdPrefix, id)
	lifememoFollowFollowCountUserIdKey := fmt.Sprintf("%s%v", cacheLifememoFollowFollowCountUserIdPrefix, userId)

	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set follow_count = follow_count + ? where user_id = ?", m.table)
		return session.ExecCtx(ctx, query, opt, userId)
	}, lifememoFollowFollowCountIdKey, lifememoFollowFollowCountUserIdKey)

	return err
}

// UpdateFansCount 更新粉丝数 1 增加 -1 减少
func (m *customFollowCountModel) UpdateFansCount(ctx context.Context, session sqlx.Session, id, userId, opt int64) error {
	lifememoFollowFollowCountIdKey := fmt.Sprintf("%s%v", cacheLifememoFollowFollowCountIdPrefix, id)
	lifememoFollowFollowCountUserIdKey := fmt.Sprintf("%s%v", cacheLifememoFollowFollowCountUserIdPrefix, userId)

	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set fans_count = fans_count + ? where user_id = ?", m.table)
		return session.ExecCtx(ctx, query, opt, userId)
	}, lifememoFollowFollowCountIdKey, lifememoFollowFollowCountUserIdKey)

	return err
}

// Trans 事务
func (m *customFollowCountModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}
