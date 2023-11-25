// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	likeRecordFieldNames          = builder.RawFieldNames(&LikeRecord{})
	likeRecordRows                = strings.Join(likeRecordFieldNames, ",")
	likeRecordRowsExpectAutoSet   = strings.Join(stringx.Remove(likeRecordFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	likeRecordRowsWithPlaceHolder = strings.Join(stringx.Remove(likeRecordFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheLifememoLikeLikeRecordIdPrefix               = "cache:lifememoLike:likeRecord:id:"
	cacheLifememoLikeLikeRecordBizIdObjIdUserIdPrefix = "cache:lifememoLike:likeRecord:bizId:objId:userId:"
)

type (
	likeRecordModel interface {
		Insert(ctx context.Context, data *LikeRecord) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*LikeRecord, error)
		FindOneByBizIdObjIdUserId(ctx context.Context, bizId string, objId int64, userId int64) (*LikeRecord, error)
		Update(ctx context.Context, data *LikeRecord) error
		Delete(ctx context.Context, id int64) error
	}

	defaultLikeRecordModel struct {
		sqlc.CachedConn
		table string
	}

	LikeRecord struct {
		Id         int64     `db:"id"`          // 主键ID
		BizId      string    `db:"biz_id"`      // 业务ID
		ObjId      int64     `db:"obj_id"`      // 点赞对象id
		UserId     int64     `db:"user_id"`     // 用户ID
		LikeType   int64     `db:"like_type"`   // 类型 0:取消点赞 1:点赞
		CreateTime time.Time `db:"create_time"` // 创建时间
		UpdateTime time.Time `db:"update_time"` // 最后修改时间
	}
)

func newLikeRecordModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultLikeRecordModel {
	return &defaultLikeRecordModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`like_record`",
	}
}

func (m *defaultLikeRecordModel) Delete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	lifememoLikeLikeRecordBizIdObjIdUserIdKey := fmt.Sprintf("%s%v:%v:%v", cacheLifememoLikeLikeRecordBizIdObjIdUserIdPrefix, data.BizId, data.ObjId, data.UserId)
	lifememoLikeLikeRecordIdKey := fmt.Sprintf("%s%v", cacheLifememoLikeLikeRecordIdPrefix, id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, lifememoLikeLikeRecordBizIdObjIdUserIdKey, lifememoLikeLikeRecordIdKey)
	return err
}

func (m *defaultLikeRecordModel) FindOne(ctx context.Context, id int64) (*LikeRecord, error) {
	lifememoLikeLikeRecordIdKey := fmt.Sprintf("%s%v", cacheLifememoLikeLikeRecordIdPrefix, id)
	var resp LikeRecord
	err := m.QueryRowCtx(ctx, &resp, lifememoLikeLikeRecordIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", likeRecordRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultLikeRecordModel) FindOneByBizIdObjIdUserId(ctx context.Context, bizId string, objId int64, userId int64) (*LikeRecord, error) {
	lifememoLikeLikeRecordBizIdObjIdUserIdKey := fmt.Sprintf("%s%v:%v:%v", cacheLifememoLikeLikeRecordBizIdObjIdUserIdPrefix, bizId, objId, userId)
	var resp LikeRecord
	err := m.QueryRowIndexCtx(ctx, &resp, lifememoLikeLikeRecordBizIdObjIdUserIdKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where `biz_id` = ? and `obj_id` = ? and `user_id` = ? limit 1", likeRecordRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, bizId, objId, userId); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultLikeRecordModel) Insert(ctx context.Context, data *LikeRecord) (sql.Result, error) {
	lifememoLikeLikeRecordBizIdObjIdUserIdKey := fmt.Sprintf("%s%v:%v:%v", cacheLifememoLikeLikeRecordBizIdObjIdUserIdPrefix, data.BizId, data.ObjId, data.UserId)
	lifememoLikeLikeRecordIdKey := fmt.Sprintf("%s%v", cacheLifememoLikeLikeRecordIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, likeRecordRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.BizId, data.ObjId, data.UserId, data.LikeType)
	}, lifememoLikeLikeRecordBizIdObjIdUserIdKey, lifememoLikeLikeRecordIdKey)
	return ret, err
}

func (m *defaultLikeRecordModel) Update(ctx context.Context, newData *LikeRecord) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	lifememoLikeLikeRecordBizIdObjIdUserIdKey := fmt.Sprintf("%s%v:%v:%v", cacheLifememoLikeLikeRecordBizIdObjIdUserIdPrefix, data.BizId, data.ObjId, data.UserId)
	lifememoLikeLikeRecordIdKey := fmt.Sprintf("%s%v", cacheLifememoLikeLikeRecordIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, likeRecordRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.BizId, newData.ObjId, newData.UserId, newData.LikeType, newData.Id)
	}, lifememoLikeLikeRecordBizIdObjIdUserIdKey, lifememoLikeLikeRecordIdKey)
	return err
}

func (m *defaultLikeRecordModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheLifememoLikeLikeRecordIdPrefix, primary)
}

func (m *defaultLikeRecordModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", likeRecordRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultLikeRecordModel) tableName() string {
	return m.table
}
