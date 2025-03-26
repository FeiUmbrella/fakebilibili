package users

import (
	"fakebilibili/adapter/http/receive"
	"fakebilibili/adapter/http/response"
	"fakebilibili/infrastructure/model/contribution/article"
	"fakebilibili/infrastructure/model/contribution/video"
	user2 "fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/model/user/attention"
	"fmt"
)

// GetAttentionList 获取关注列表
func GetAttentionList(data *receive.GetAttentionListReceiveStruct, uid uint) (interface{}, error) {
	// 获取某个用户的关注列表
	somebodyAtnList := new(attention.AttentionsList)
	err := somebodyAtnList.GetAttentionList(data.ID)
	if err != nil {
		return nil, fmt.Errorf("获取查询用户关注列表失败")
	}
	// 获取自己关注的用户的id
	selfAL := new(attention.AttentionsList)
	attentionIdList, err := selfAL.GetAttentionListByIdArr(uid)
	if err != nil {
		return nil, fmt.Errorf("获取自己关注列表ID失败")
	}
	res, err := response.GetAttentionListResponse(somebodyAtnList, attentionIdList)
	if err != nil {
		return nil, fmt.Errorf("响应失败")
	}
	return res, nil
}

// GetVermicelliList 获取粉丝列表
func GetVermicelliList(data *receive.GetVermicelliListReceiveStruct, uid uint) (interface{}, error) {
	// 获取某个用户的粉丝列表
	somebodyAtnList := new(attention.AttentionsList)
	err := somebodyAtnList.GetVermicelliList(data.ID)
	if err != nil {
		return nil, fmt.Errorf("获取查询用户粉丝列表失败")
	}
	// 获取自己的粉丝id
	selfAL := new(attention.AttentionsList)
	attentionIdList, err := selfAL.GetAttentionListByIdArr(uid)
	if err != nil {
		return nil, fmt.Errorf("获取自己粉丝列表ID失败")
	}
	res, err := response.GetVermicelliListResponse(somebodyAtnList, attentionIdList)
	if err != nil {
		return nil, fmt.Errorf("响应失败")
	}
	return res, nil
}

// GetSpaceIndividual 获取个人空间
func GetSpaceIndividual(data *receive.GetSpaceIndividualReceiveStruct, uid uint) (interface{}, error) {
	// 获取用户信息
	user := new(user2.User)
	user.Find(data.ID)
	isAttention := false
	at := new(attention.Attention)
	if uid != 0 {
		// 查询 uid(用户自己) 是否关注 data.ID
		isAttention = at.IsAttention(uid, data.ID)
	}
	attentionNum, err := at.GetAttentionNum(data.ID)
	if err != nil {
		return nil, fmt.Errorf("查询 %s 的关注数量失败", user.Username)
	}
	vermicelliNum, err := at.GetVermicelliNum(data.ID)
	if err != nil {
		return nil, fmt.Errorf("查询 %s 的粉丝数量失败", user.Username)
	}
	return response.GetSpaceIndividualResponse(user, isAttention, attentionNum, vermicelliNum)
}

// GetReleaseInformation 获取用户作品(视频and专栏)
func GetReleaseInformation(data *receive.GetReleaseInformationReceiveStruct) (interface{}, error) {
	// 获取视频列表
	videoList := new(video.VideosContributionList)
	err := videoList.GetVideoListBySpace(data.ID)
	if err != nil {
		return nil, fmt.Errorf("查询空间视频失败")
	}
	// 获取专栏列表
	articleList := new(article.ArticlesContributionList)
	err = articleList.GetArticleBySpace(data.ID)
	if err != nil {
		return nil, fmt.Errorf("查询空间专栏失败")
	}
	res, err := response.GetReleaseInformationResponse(videoList, articleList)
	if err != nil {
		return nil, fmt.Errorf("响应失败")
	}
	return res, nil
}
