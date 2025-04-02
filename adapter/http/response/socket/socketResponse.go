package socket

type NoticeMessageStruct struct {
	NoticeType string `json:"notice_type"`
	Unread     *int64 `json:"unread"` // 未读通知的数量
}
