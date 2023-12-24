package types

const (
	FollowStatusFollow   = iota + 1 // 关注
	FollowStatusUnfollow            // 取消关注
)

const (
	DefaultLimit        = 20   // 默认分页大小
	CacheMaxFollowCount = 1000 // 缓存最大关注数
	CacheMaxFansCount   = 1000 // 缓存最大粉丝数
)
