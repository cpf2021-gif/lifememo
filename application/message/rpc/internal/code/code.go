package code

import "lifememo/pkg/xcode"

var (
	NoSendToSelf   = xcode.New(40000, "不能给自己发消息")
	NoPermission   = xcode.New(40001, "没有权限")
	DeletedMessage = xcode.New(40002, "消息已删除")
	MessageInvalid = xcode.New(40003, "消息不存在")
	NoneMessage    = xcode.New(40004, "没有消息")
)
