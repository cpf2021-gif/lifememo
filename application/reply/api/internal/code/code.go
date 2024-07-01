package code

import "lifememo/pkg/xcode"

var (
	ReplyNotFound         = xcode.New(50000, "评论不存在")
	ErrInvalidBizId       = xcode.New(50001, "无效的业务ID")
	ErrInvalidTargetId    = xcode.New(50002, "无效的目标ID")
	ErrInvalidParentId    = xcode.New(50003, "无效的父级ID")
	ErrInvalidBeRepliedId = xcode.New(50004, "无效的被回复者ID")
	BeRepliedUserNotFound = xcode.New(50005, "被回复者不存在")
	ReplyIsNotRoot        = xcode.New(50006, "该评论不是根评论")
)
