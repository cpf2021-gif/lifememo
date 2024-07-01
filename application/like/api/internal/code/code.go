package code

import "lifememo/pkg/xcode"

var (
	BizIdError       = xcode.New(30000, "业务错误")
	ObjIdInvalid     = xcode.New(30001, "对象ID无效")
	TypeInvalid      = xcode.New(30002, "点赞类型无效")
	NotThumbuped     = xcode.New(30003, "未点赞")
	AlreadyThumbuped = xcode.New(30004, "已点赞")
)
