package home

import "fakebilibili/infrastructure/model/common"

// GetHomeInfoReceiveStruct 用户访问主页（包含轮播图什么的）需要传入的参数
type GetHomeInfoReceiveStruct struct {
	PageInfo common.PageInfo `json:"page_info" binding:"required"`
}

// SubmitBugReceiveStruct 用户反馈bug
type SubmitBugReceiveStruct struct {
	Content string `json:"content" binding:"required" form:"content"`
	Phone   string `json:"phone" form:"phone"`
}
