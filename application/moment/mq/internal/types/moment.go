package types

// CanalMsg 收集canal解析binlog后的数据
type CanalMsg struct {
	Data []struct {
		ID          string `json:"id"`
		Content     string `json:"content"`
		AuthorID    string `json:"author_id"`
		Status      string `json:"status"`
		CommentNum  string `json:"comment_num"`
		LikeNum     string `json:"like_num"`
		CollectNum  string `json:"collect_num"`
		ViewNum     string `json:"view_num"`
		ShareNum    string `json:"share_num"`
		TagIds      string `json:"tag_ids"`
		PublishTime string `json:"publish_time"`
		CreateTime  string `json:"create_time"`
		UpdateTime  string `json:"update_time"`
	} `json:"data"`

	Type string `json:"type"`
}

type MomentEsMsg struct {
	MomentID    int64  `json:"moment_id"`
	Content     string `json:"content"`
	AuthorID    int64  `json:"author_id"`
	Status      int    `json:"status"`
	CommentNum  int64  `json:"comment_num"`
	LikeNum     int64  `json:"like_num"`
	CollectNum  int64  `json:"collect_num"`
	ViewNum     int64  `json:"view_num"`
	ShareNum    int64  `json:"share_num"`
	TagIds      string `json:"tag_ids"`
	PublishTime string `json:"publish_time"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
}
