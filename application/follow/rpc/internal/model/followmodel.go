package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FollowModel = (*customFollowModel)(nil)

type (
	// FollowModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFollowModel.
	FollowModel interface {
		followModel
		FollowListByUserId(ctx context.Context, userId, limit int64, cursor string) ([]*Follow, error)
		FansListByUserId(ctx context.Context, userId, limit int64, cursor string) ([]*Follow, error)
		InsertWithSession(ctx context.Context, session sqlx.Session, data *Follow) (sql.Result, error)
		UpdateWithSession(ctx context.Context, session sqlx.Session, data *Follow) error
	}

	customFollowModel struct {
		*defaultFollowModel
	}
)

// NewFollowModel returns a model for the database table.
func NewFollowModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FollowModel {
	return &customFollowModel{
		defaultFollowModel: newFollowModel(conn, c, opts...),
	}
}

func (m *customFollowModel) FollowListByUserId(ctx context.Context, userId, limit int64, cursor string) ([]*Follow, error) {
	var (
		err     error
		query   string
		follows []*Follow
	)
	query = "select " + followRows + " from " + m.table + " where user_id=? and follow_status=1 and create_time<=? order by id desc limit ?"
	err = m.QueryRowsNoCacheCtx(ctx, &follows, query, userId, cursor, limit)
	if err != nil {
		return nil, err
	}

	return follows, nil
}

func (m *customFollowModel) FansListByUserId(ctx context.Context, userId, limit int64, cursor string) ([]*Follow, error) {
	var (
		err     error
		query   string
		follows []*Follow
	)
	query = "select " + followRows + " from " + m.table + " where followed_user_id=? and follow_status=1 and create_time<=? order by id desc limit ?"
	err = m.QueryRowsNoCacheCtx(ctx, &follows, query, userId, cursor, limit)
	if err != nil {
		return nil, err
	}

	return follows, nil
}

func (m *customFollowModel) InsertWithSession(ctx context.Context, session sqlx.Session, data *Follow) (sql.Result, error) {
	lifememoFollowFollowIdKey := fmt.Sprintf("%s%v", cacheLifememoFollowFollowIdPrefix, data.Id)
	lifememoFollowFollowUserIdFollowedUserIdKey := fmt.Sprintf("%s%v:%v", cacheLifememoFollowFollowUserIdFollowedUserIdPrefix, data.UserId, data.FollowedUserId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, followRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.UserId, data.FollowedUserId, data.FollowStatus)
	}, lifememoFollowFollowIdKey, lifememoFollowFollowUserIdFollowedUserIdKey)
	return ret, err
}

func (m *customFollowModel) UpdateWithSession(ctx context.Context, session sqlx.Session, newData *Follow) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	lifememoFollowFollowIdKey := fmt.Sprintf("%s%v", cacheLifememoFollowFollowIdPrefix, data.Id)
	lifememoFollowFollowUserIdFollowedUserIdKey := fmt.Sprintf("%s%v:%v", cacheLifememoFollowFollowUserIdFollowedUserIdPrefix, data.UserId, data.FollowedUserId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, followRowsWithPlaceHolder)
		return session.ExecCtx(ctx, query, newData.UserId, newData.FollowedUserId, newData.FollowStatus, newData.Id)
	}, lifememoFollowFollowIdKey, lifememoFollowFollowUserIdFollowedUserIdKey)
	return err
}
