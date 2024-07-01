package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MessageModel = (*customMessageModel)(nil)

type (
	// MessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMessageModel.
	MessageModel interface {
		messageModel
		FindMessageBySenderAndReceiver(sender int64, receiver int64) ([]*Message, error)
	}

	customMessageModel struct {
		*defaultMessageModel
	}
)

// NewMessageModel returns a model for the database table.
func NewMessageModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MessageModel {
	return &customMessageModel{
		defaultMessageModel: newMessageModel(conn, c, opts...),
	}
}

func (m *customMessageModel) FindMessageBySenderAndReceiver(sender int64, receiver int64) ([]*Message, error) {
	var (
		err      error
		query    string
		messages []*Message
	)

	query = "select " + messageRows + " from " + m.table + " where (send_user_id=? and receive_user_id=? or send_user_id=? and receive_user_id=?) and status=1 order by create_time"

	err = m.QueryRowsNoCache(&messages, query, sender, receiver, receiver, sender)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
