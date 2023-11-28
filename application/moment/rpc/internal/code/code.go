package code

import "lifememo/pkg/xcode"

var (
	SortTypeInvalid     = xcode.New(20000, "排序类型无效")
	MomentContentEmpty  = xcode.New(20001, "动态内容不能为空")
	MomentIdInvalid     = xcode.New(20002, "动态ID无效")
	UserIdInvalid       = xcode.New(20003, "用户ID无效")
	MomentStatusInvalid = xcode.New(20004, "动态状态无效")
)
