package oss

import (
	"errors"
	"fakebilibili/infrastructure/pkg/global"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ice20201109 "github.com/alibabacloud-go/ice-20201109/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
	"strings"
)

var (
	accessKeyId     = global.Config.AliyunOss.AccessKeyId
	accessKeySecret = global.Config.AliyunOss.AccessKeySecret
	roleArn         = global.Config.AliyunOss.RoleArn
	roleSessionName = global.Config.AliyunOss.RoleSessionName
	durationSeconds = global.Config.AliyunOss.DurationSeconds
	endpoint        = global.Config.AliyunOss.Endpoint
	ossEndpoint     = global.Config.AliyunOss.OssEndPoint
	bucketName      = global.Config.AliyunOss.Bucket
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

// GetMediaInfo 利用将oss视频文件注册媒体资源得到的ID，获取对应媒体资源信息
func GetMediaInfo(mediaID *string) (*ice20201109.GetMediaInfoResponse, error) {
	client, err := CreateIceClient(tea.String(accessKeyId), tea.String(accessKeySecret))
	if err != nil {
		global.Logger.Errorf("初始化Ice_cilent失败 err : %s", err.Error())
	}
	getMediaInfoRequest := &ice20201109.GetMediaInfoRequest{
		MediaId: mediaID,
	}
	runtime := &util.RuntimeOptions{}
	if r := tea.Recover(recover()); r != nil {
	}

	res, err := client.GetMediaInfoWithOptions(getMediaInfoRequest, runtime)
	global.Logger.Infof("打印iceclient返回的GetMediaInfo的结果：%v", res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SubmitTranscodeJob 云转码向媒资提交转码任务
func SubmitTranscodeJob(taskName, mediaID, outputUrl, template string) (*ice20201109.SubmitTranscodeJobResponseBody, error) {
	client, err := CreateIceClient(tea.String(accessKeyId), tea.String(accessKeySecret))
	if err != nil {
		global.Logger.Errorf("初始化Ice_cilent失败 err : %s", err.Error())
	}
	inputGroup0 := &ice20201109.SubmitTranscodeJobRequestInputGroup{
		Type:  tea.String("Media"),
		Media: tea.String(mediaID),
	}
	outputGroup0Output := &ice20201109.SubmitTranscodeJobRequestOutputGroupOutput{
		Type:  tea.String("OSS"),
		Media: tea.String(outputUrl),
	}
	outputGroup0ProcessConfigTranscode := &ice20201109.SubmitTranscodeJobRequestOutputGroupProcessConfigTranscode{
		TemplateId: tea.String(template),
	}
	outputGroup0ProcessConfig := &ice20201109.SubmitTranscodeJobRequestOutputGroupProcessConfig{
		Transcode: outputGroup0ProcessConfigTranscode,
	}
	outputGroup0 := &ice20201109.SubmitTranscodeJobRequestOutputGroup{
		ProcessConfig: outputGroup0ProcessConfig,
		Output:        outputGroup0Output,
	}
	submitTranscodeJobRequest := &ice20201109.SubmitTranscodeJobRequest{
		OutputGroup: []*ice20201109.SubmitTranscodeJobRequestOutputGroup{outputGroup0},
		Name:        tea.String(taskName),
		InputGroup:  []*ice20201109.SubmitTranscodeJobRequestInputGroup{inputGroup0},
	}
	runtime := &util.RuntimeOptions{}
	defer func() {
		if r := tea.Recover(recover()); r != nil {
		}
	}()
	result, err := client.SubmitTranscodeJobWithOptions(submitTranscodeJobRequest, runtime)
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

// CreateOSSClient 获取oss文件操作对象
func CreateOSSClient() (*oss.Client, error) {
	_ = os.Setenv("OSS_ACCESS_KEY_ID", accessKeyId)
	_ = os.Setenv("OSS_ACCESS_KEY_SECRET", accessKeySecret)
	//这个方法是从环境变量里读取id和key来生成authorization信息的,所以需要先设置环境变量
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()

	client, err := oss.New(ossEndpoint, accessKeyId, accessKeySecret, oss.SetCredentialsProvider(&provider))
	if err != nil {
		return nil, err
	}
	return client, nil
}

// DeleteOSSFile 删除oss中的视频
func DeleteOSSFile(filePath []string) error {
	client, err := CreateOSSClient()
	if err != nil {
		return errors.New("创建 oss_client err: " + err.Error())
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return errors.New("获取bucket err : " + err.Error())
	}

	// 删除给定文件
	for _, path := range filePath {
		// 预处理一下path，因为path除了实际路径还包括src: 这个前缀，以及type:oss这个后缀
		result := strings.Split(path, `"src": "`)
		result = strings.Split(result[1], `", "type": "oss"`)
		//此时result只剩一个元素，也是我们需要的最终路径
		endPath := result[0]

		global.Logger.Infof("要删除的路径为%s", endPath)
		err := bucket.DeleteObject(endPath)
		if err != nil {
			return errors.New("删除文件 err : " + err.Error())
		}
	}
	return nil
}
