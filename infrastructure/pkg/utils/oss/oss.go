package oss

import (
	"fakebilibili/infrastructure/pkg/global"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ice20201109 "github.com/alibabacloud-go/ice-20201109/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

var (
	accessKeyId     = global.Config.AliyunOss.AccessKeyId
	accessKeySecret = global.Config.AliyunOss.AccessKeySecret
	roleArn         = global.Config.AliyunOss.RoleArn
	roleSessionName = global.Config.AliyunOss.RoleSessionName
	durationSeconds = global.Config.AliyunOss.DurationSeconds
	endpoint        = global.Config.AliyunOss.Endpoint
)

// CreateStsClient 创建STS Client
func CreateStsClient(accessKeyId *string, accessKeySecret *string) (*sts20150401.Client, error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String(endpoint)
	res, err := sts20150401.NewClient(config)
	return res, err
}

// GetStsInfo 获取STS临时授权码
func GetStsInfo() (*sts20150401.AssumeRoleResponseBodyCredentials, error) {
	client, err := CreateStsClient(tea.String(accessKeyId), tea.String(accessKeySecret))
	if err != nil {
		return nil, fmt.Errorf("CreateStsClient方法创建client出错: %s", err.Error())
	}
	assumeRoleRequest := &sts20150401.AssumeRoleRequest{
		RoleArn:         tea.String(roleArn),
		RoleSessionName: tea.String(roleSessionName),
		//DurationSeconds 申请临时授权码有效时间 最小15分钟，最大一个小时
		DurationSeconds: tea.Int64(3600),
	}
	runtime := &util.RuntimeOptions{} // 设置超时
	defer func() {
		if r := tea.Recover(recover()); r != nil {
		}
	}()
	res, err := client.AssumeRoleWithOptions(assumeRoleRequest, runtime)
	if err != nil {
		return nil, err
	}
	if *res.StatusCode != 200 {
		return nil, fmt.Errorf("错误的状态码: %d", res.StatusCode)
	}
	return res.Body.Credentials, nil
}

// CreateIceClient 创建注册媒资的Client
func CreateIceClient(accessKeyId *string, accessKeySecret *string) (*ice20201109.Client, error) {
	config := &openapi.Config{
		RegionId:        tea.String("cn-hangzhou"),
		Endpoint:        tea.String("ice.cn-hangzhou.aliyuncs.com"),
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	client, err := ice20201109.NewClient(config)
	return client, err
}

// RegisterMediaInfo 注册媒体资源
// params:视频的存储路径(https://<bucket_name>.<region>.aliyuncs.com/E:/video/hash.mp4)；媒体资源的类型；和当前时间转成的字符串
func RegisterMediaInfo(inputUrl, mediaType, Title string) (*ice20201109.RegisterMediaInfoResponseBody, error) {
	client, err := CreateIceClient(tea.String(accessKeyId), tea.String(accessKeySecret))
	if err != nil {
		global.Logger.Errorf("初始化媒资cilent失败 err : %s", err.Error())
	}

	registerMediaInfoRequest := &ice20201109.RegisterMediaInfoRequest{
		Overwrite: tea.Bool(true),        // 是否覆盖已注册媒资，默认 false。-true，如果 inputUrl 已注册，删除原有媒资并注册新媒资；
		InputURL:  tea.String(inputUrl),  // 资源在OSS的位置
		MediaType: tea.String(mediaType), // 资源类型
		Title:     tea.String(Title),     // 标题，若不提供，根据日期自动生成默认 title。
	}
	res, err := client.RegisterMediaInfo(registerMediaInfoRequest)
	if err != nil {
		global.Logger.Errorf("注册媒资失败 err %s ", err.Error())
	}
	if *res.StatusCode != 200 {
		global.Logger.Errorf("注册媒体资源失败，返回的错误信息为:" + err.Error())
		return nil, err
	}
	return res.Body, nil
}
