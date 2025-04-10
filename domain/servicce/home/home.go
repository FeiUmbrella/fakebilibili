package home

import (
	"fakebilibili/adapter/http/receive/home"
	home3 "fakebilibili/adapter/http/response/home"
	"fakebilibili/infrastructure/consts"
	"fakebilibili/infrastructure/model/contribution/video"
	home2 "fakebilibili/infrastructure/model/home"
	"fakebilibili/infrastructure/pkg/global"
	"fakebilibili/infrastructure/pkg/utils/email"
	"fakebilibili/infrastructure/pkg/utils/validator"
	"fmt"
)

// GetHomeInfo 获取主页轮播图和推荐视频
func GetHomeInfo(data *home.GetHomeInfoReceiveStruct) (interface{}, error) {
	// 获取主页轮播图
	rotoGraphList := new(home2.List)
	err := rotoGraphList.GetALL()
	if err != nil {
		return nil, err
	}

	// 获取主页推荐视频
	videoList := new(video.VideosContributionList)
	err = videoList.GetHomeVideoList(data.PageInfo)
	if err != nil {
		return nil, err
	}
	res := &home3.GetHomeInfoResponse{}
	res.Response(rotoGraphList, videoList)
	return res, nil
}

// SubmitBug 用户反馈bug
func SubmitBug(data *home.SubmitBugReceiveStruct) (interface{}, error) {
	// 给小号发个邮件
	emailTo := []string{consts.SystemEmail}
	err := email.SendEmail(emailTo, "用户反馈的bug信息", fmt.Sprintf("用户反馈的bug信息为：%s\n用户留下的联系方式为:%s", data.Content, data.Phone))
	if err != nil {
		global.Logger.Errorf("向系统邮箱发送bug反馈信息失败：%v", err)
		return "保存bug信息失败", err
	}
	phone := data.Phone
	if validator.VerifyMobileFormat(phone) {
		// todo: 这是回短信通知用户
	}
	if validator.VerifyEmailFormat(phone) {
		//可以给人回个邮件答复一下
		if err = email.SendEmail([]string{data.Phone}, "感谢您的反馈", "已收到您反馈的问题，处理完成后会给您答复"); err != nil {
			global.Logger.Errorf("给用户反馈邮件失败:%v", err)
		}
		global.Logger.Infoln("成功给用户推送确认信息邮件")
	}
	return "success", nil
}
