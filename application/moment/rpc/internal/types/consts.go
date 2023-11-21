package types

const (
	SortPublishTime = iota
	SortLikeCount
)

const (
	DefaultPageSize = 20
	DefaultLimit    = 200

	DefaultSortLikeCursor = 1 << 30
)

const (
	// MomentStatusPending 待审核
	MomentStatusPending = iota
	// MomentStatusNotPass 审核不通过
	MomentStatusNotPass
	// MomentStatusVisible 可见
	MomentStatusVisible
	// MomentStatusPrivate 私密
	MomentStatusPrivate
	// MomentStatusUserDelete 用户删除
	MomentStatusUserDelete
)
