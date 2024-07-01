package types

type ThumbupMsg struct {
	BizId    string ` json:"bizId,omitempty"`    // 业务id
	ObjId    int64  ` json:"objId,omitempty"`    // 点赞对象id
	UserId   int64  ` json:"userId,omitempty"`   // 用户id
	LikeType int64  ` json:"likeType,omitempty"` // 点赞类型   0:取消点赞  1:点赞
}
