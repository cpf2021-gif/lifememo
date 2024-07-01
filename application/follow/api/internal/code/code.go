package code

import "lifememo/pkg/xcode"

var (
	FollowUserIdInvalid = xcode.New(60000, "关注用户ID无效")
	NoFollowed          = xcode.New(60001, "未关注该用户")
	DontFollowYourself  = xcode.New(60002, "不能对自己进行操作")
	Followed            = xcode.New(60003, "已关注该用户")
)
