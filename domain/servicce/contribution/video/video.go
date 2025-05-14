package video

import (
	"encoding/json"
	"errors"
	"fakebilibili/adapter/http/receive/contribution/video"
	video2 "fakebilibili/adapter/http/response/contribution/video"
	"fakebilibili/domain/servicce/contribution/socket"
	"fakebilibili/domain/servicce/users"
	"fakebilibili/infrastructure/consts"
	"fakebilibili/infrastructure/model/common"
	video3 "fakebilibili/infrastructure/model/contribution/video"
	"fakebilibili/infrastructure/model/contribution/video/barrage"
	"fakebilibili/infrastructure/model/contribution/video/comments"
	"fakebilibili/infrastructure/model/sundry"
	"fakebilibili/infrastructure/model/user/attention"
	"fakebilibili/infrastructure/model/user/collect"
	"fakebilibili/infrastructure/model/user/favorites"
	"fakebilibili/infrastructure/model/user/notice"
	"fakebilibili/infrastructure/model/user/record"
	"fakebilibili/infrastructure/pkg/global"
	"fakebilibili/infrastructure/pkg/utils/calculate"
	"fakebilibili/infrastructure/pkg/utils/conversion"
	"fakebilibili/infrastructure/pkg/utils/oss"
	"fmt"
	"github.com/FeiUmbrella/Sensitive_Words_Filter/filter"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// GetVideoBarrage  获取视频弹幕 (播放器）
func GetVideoBarrage(data *video.GetVideoBarrageReceiveStruct) (interface{}, error) {
	list := new(barrage.BarragesList)
	videoID, _ := strconv.ParseUint(data.ID, 0, 19) // base=0 根据参数前缀判断进制
	if !list.GetVideoBarrageByID(uint(videoID)) {
		return nil, fmt.Errorf("查询视频弹幕失败")
	}
	res := video2.GetVideoBarrageResponse(list)
	return res, nil
}

// GetVideoBarrageList 获取视频弹幕，先从redis中获取
func GetVideoBarrageList(data *video.GetVideoBarrageListReceiveStruct) (results interface{}, err error) {
	list := new(barrage.BarragesList)

	// 查redis缓存
	key := fmt.Sprintf("%s%s", consts.VideoBarragePrefix, data.ID)
	res, err := global.RedisDb.Get(key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		global.Logger.Errorf("获取视频%v弹幕信息时，查询redis err:%v", data.ID, err)
	}
	if len(res) != 0 {
		// 查到了弹幕缓存，直接封装返回
		if err := json.Unmarshal([]byte(res), list); err == nil {
			global.Logger.Infof("查询视频%s弹幕信息时，命中cache返回：%v", data.ID, list)
			return list, nil
		}
		global.Logger.Errorf("获取视频%v弹幕信息时，解封装错误", data.ID)
	}
	// 数据库中查找
	videoID, _ := strconv.ParseUint(data.ID, 0, 19) // base=0 根据参数前缀判断进制
	if !list.GetVideoBarrageByID(uint(videoID)) {
		return nil, fmt.Errorf("查询视频弹幕失败")
	}

	// set redis
	bytes, err := json.Marshal(list)
	if err != nil {
		global.Logger.Errorf("封装视频%v弹幕信息错误:%v", data.ID, err)
	} else {
		global.RedisDb.Set(key, bytes, 5*time.Second) // 缓存5s
	}
	results = video2.GetVideoBarrageListResponse(list)
	return results, nil
}

// GetVideoComment 获取视频评论
func GetVideoComment(data *video.GetVideoCommentReceiveStruct) (results interface{}, err error) {
	fmt.Println(data.VideoID, data.PageInfo)
	videosContribution := new(video3.VideosContribution)
	if !videosContribution.GetVideoComments(data.VideoID, data.PageInfo) {
		return nil, fmt.Errorf("查询失败")
	}
	return video2.GetVideoContributionCommentsResponse(videosContribution), nil
}

// LikeVideoComment 给视频点赞
func LikeVideoComment(data *video.LikeVideoCommentReqStruct) (results interface{}, err error) {
	zsetKey := fmt.Sprintf("%s%s", consts.VideoCommentZSetPrefix, strconv.Itoa(data.VideoCommentId))
	hashKey := fmt.Sprintf("%s%s", consts.VideoCommentHashPrefix, strconv.Itoa(data.VideoCommentId))

	// 判断这条评论是否在Hash中
	_, err = global.RedisDb.HGet(hashKey, strconv.Itoa(data.VideoCommentId)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		global.Logger.Errorf("在hash里查询评论点赞数failed:%v", err)
	}
	if errors.Is(err, redis.Nil) { // 这条评论未放进Hash中
		global.Logger.Infof("向hash和zset里添加id=%d的评论", data.VideoCommentId)
		com := &comments.Comment{}
		if err := global.MysqlDb.Model(&comments.Comment{}).Where("id = ?", data.VideoCommentId).First(&com).Error; err != nil {
			global.Logger.Errorf("从数据库查询id=%d的评论内容失败: %v", data.VideoCommentId, err)
		}
		com.Heat++ // 点赞
		// 将评论放入Hash
		_, err := global.RedisDb.HSet(hashKey, strconv.Itoa(data.VideoCommentId), com).Result()
		if err != nil {
			global.Logger.Errorf("向hash中放入id=%d的评论内容失败: %v", data.VideoCommentId, err)
		}
		// 将评论放入 Zset
		jsonComment, _ := json.Marshal(com)
		_, err = global.RedisDb.ZAdd(zsetKey, redis.Z{
			Member: jsonComment,
			Score:  float64(com.Heat),
		}).Result()
		if err != nil {
			global.Logger.Errorf("在zset里增加评论 %v failed:%v", com, err)
		}
	} else {
		global.Logger.Infof("更新id=%d的评论", data.VideoCommentId)
		// 该条评论已经在Hash中，从Hash中取出评论，作为ZIncrBy参数实现heat+1
		res, err := global.RedisDb.HGet(hashKey, strconv.Itoa(data.VideoCommentId)).Result()
		if err != nil {
			global.Logger.Errorf("从hash中获取id=%d的评论failed:%v", data.VideoCommentId, err)
		}
		_, err = global.RedisDb.ZIncrBy(zsetKey, 1, res).Result()
		if err != nil {
			global.Logger.Errorf("向zset中评论 ： %v 的热度+1failed: %v", res, err)
		}
	}
	return 1, nil
}

// GetVideoContributionByID 获取视频信息
func GetVideoContributionByID(data *video.GetVideoContributionByIDReceiveStruct, uid uint) (results interface{}, err error) {
	// 观看视频同时将视频id放入bitmap，推荐视频时随即请求，然后过滤掉最近推荐过的和观看过的
	key := fmt.Sprintf("%s%d", consts.UniqueVideoRecommendPrefix, -1)
	_, err = global.RedisDb.SetBit(key, int64(data.VideoID), 1).Result()
	if err != nil {
		global.Logger.Errorf("set bitmap failed:%v", err)
	}
	videoInfo := new(video3.VideosContribution)
	err = videoInfo.FindByID(data.VideoID)
	if err != nil {
		return nil, fmt.Errorf("查询视频信息失败")
	}
	isAttention := false
	isLike := false
	isCollect := false
	if uid != 0 {
		// 当被登录用户点击视频时才增加播放量
		//todo:这里有个bug：当redis缓存里有这个视频的信息时，点击视频不增加播放量
		//bug的原因是加载播放页面时useInit函数和vidoe/video.vue的onMounted各调用了一次GetVideoContributionByID方法，就导致heat递增2次
		if !global.RedisDb.SIsMember(consts.VideoWatchByID+strconv.Itoa(int(data.VideoID)), uid).Val() {
			//SIsMember key value :查询redis set的集合key中是否存在value
			// 最近无播放，将用户添加在这个视频的set中
			global.RedisDb.SAdd(consts.VideoWatchByID+strconv.Itoa(int(data.VideoID)), uid)
			if videoInfo.Watch(data.VideoID) != nil {
				// 更新数据库
				global.Logger.Error("添加播放量错误,视频video_id:", data.VideoID)
			}
			videoInfo.Heat++ // 这里更新的是从数据库取出的旧值
		}
		// 是否关注视频作者
		at := new(attention.Attention)
		isAttention = at.IsAttention(uid, videoInfo.UserInfo.ID)

		// 是否点赞该视频
		lk := new(video3.Likes)
		isLike = lk.IsLike(uid, videoInfo.ID)

		// 是否收藏
		fl := new(favorites.FavoriteList)
		err = fl.GetFavoritesList(uid) // 获取播放者的收藏夹
		if err != nil {
			return nil, fmt.Errorf("查询播放者收藏夹失败")
		}
		flIDs := make([]uint, 0)
		for _, v := range *fl {
			flIDs = append(flIDs, v.ID)
		}
		cl := new(collect.CollectsList)
		isCollect = cl.FindIsCollectByFavorites(data.VideoID, flIDs)

		// 添加历史记录
		rd := new(record.Record)
		err = rd.AddVideoRecord(uid, data.VideoID)
		if err != nil {
			return nil, fmt.Errorf("添加历史记录失败")
		}
	}
	// 获取推荐列表
	recommendList := new(video3.VideosContributionList)
	err = recommendList.GetRecommendList(uid)
	if err != nil {
		return nil, err
	}
	res := video2.GetVideoContributionByIDResponse(videoInfo, recommendList, isAttention, isLike, isCollect)
	return res, nil
}

// SendVideoBarrage 发送视频弹幕
// 获取弹幕list的时候先查了缓存，根据cache aside，这里应该先修改数据库后删除缓存
func SendVideoBarrage(data *video.SendVideoBarrageReceiveStruct, uid uint) (results interface{}, err error) {
	node := filter.NewNode()
	node.StartMatch(data.Text)
	if node.IsSensitive() {
		return nil, fmt.Errorf("弹幕中包含敏感词")
	}
	// 保存弹幕
	videoID, _ := strconv.ParseUint(data.ID, 0, 19)
	bg := barrage.Barrage{
		Uid:     uid,
		VideoID: uint(videoID),
		Time:    data.Time,
		Author:  data.Author,
		Type:    data.Type,
		Text:    data.Text,
		Color:   data.Color,
	}
	if !bg.Create() {
		return data, fmt.Errorf("发送视频弹幕失败")
	}
	// 删除缓存
	key := fmt.Sprintf("%s%s", consts.VideoBarragePrefix, data.ID)
	_ = global.RedisDb.Del(key)

	// socket 弹幕通知
	res := socket.MsgInfo{
		Type: consts.VideoSocketTypeResponseBarrageNum,
		Data: nil,
	}
	for _, v := range socket.Severe.VideoRoom[uint(videoID)] {
		v.MegList <- res
	}
	return data, nil
}

// CreateVideoContribution 创建视频信息
func CreateVideoContribution(data *video.CreateVideoContributionReceiveStruct, uid uint) (results interface{}, err error) {
	//1.处理视频在数据库中的基本信息
	videoSrc, _ := json.Marshal(common.Img{
		Src: data.VideoPath,
		Tp:  data.VideoUploadType, // local/aliyunOss
	})
	coverImg, _ := json.Marshal(common.Img{
		Src: data.CoverPath,
		Tp:  data.CoverUploadType,
	})
	var width, height int // 视频的宽，高
	// 本地上传基本启用，都上传到oss中
	if data.VideoUploadType == "local" { // 文件已上传到本地
		fmt.Println("传入 ffprobe 的路径：", data.VideoPath)
		width, height, err = calculate.GetVideoResolution(data.VideoPath)
		fmt.Println("传入视频的高宽 w:", width, "h:", height)
		if err != nil {
			global.Logger.Error("获取视频分辨率失败")
			return nil, fmt.Errorf("获取视频分辨率失败")
		}
	} else { // aliyunOss
		mediaInfo, err := oss.GetMediaInfo(data.Media)
		if err != nil {
			return nil, errors.New("获取视频信息失败，错误：" + err.Error())
		}
		width, _ = strconv.Atoi(*mediaInfo.Body.MediaInfo.FileInfoList[0].FileBasicInfo.Width)
		height, _ = strconv.Atoi(*mediaInfo.Body.MediaInfo.FileInfoList[0].FileBasicInfo.Height)
	}
	videoContribution := &video3.VideosContribution{
		// 将视频信息存放在数据库中
		Uid:           uid,
		Title:         data.Title,
		Cover:         coverImg,
		Reprinted:     conversion.Bool2Int8(*data.Reprinted),
		Label:         strings.Join(data.Label, ","),
		VideoDuration: data.VideoDuration,
		Introduce:     data.Introduce,
		Heat:          0, // 新视频播放量为0
	}
	if data.Media != nil { // 说明视频存放在Oss
		videoContribution.MediaID = *data.Media
	}

	// 2.定义分辨率列表，将每个分辨率的视频源都设为初始源
	// 高分辨率能转向低分辨率 1080p表示像素有1080行
	resolution := []int{1080, 720, 480, 360}
	if height >= 1080 {
		resolution = resolution[1:]
		videoContribution.Video = videoSrc
	} else if height >= 720 {
		resolution = resolution[2:]
		videoContribution.Video720p = videoSrc
	} else if height >= 480 {
		resolution = resolution[3:]
		videoContribution.Video480p = videoSrc
	} else if height >= 360 {
		resolution = resolution[4:]
		videoContribution.Video360p = videoSrc
	} else {
		global.Logger.Error("上传视频分辨率过低")
		return nil, fmt.Errorf("上传视频分辨率过低")
	}

	// 当没有传递定时的发布时间时，默认is_visible=1
	if data.DateTime == "" {
		videoContribution.IsVisible = 1
	}
	if !videoContribution.Create() {
		return nil, fmt.Errorf("视频信息保存失败")
	}

	// 3.进行视频转码
	var wg sync.WaitGroup
	wg.Add(1)
	go func(width, height int, video *video3.VideosContribution) {
		defer wg.Done()
		// 3.1 本地视频使用ffmpeg转码
		if data.VideoUploadType == "local" {
			inputFile := data.VideoPath
			sr := strings.Split(inputFile, ".")
			// 对于每个分辨率，将视频转换为对应分辨率并保存src
			for _, r := range resolution {
				// 计算转码后的宽，高。要等比例缩小
				w := int(math.Ceil(float64(r) / float64(height) * float64(width)))
				h := int(math.Ceil(float64(r)))
				if h >= height {
					continue
				}

				dest := sr[0] + fmt.Sprintf("_output_%dp."+sr[1], r)
				cmd := exec.Command("ffmpeg",
					"-i",
					inputFile,
					"-vf",
					fmt.Sprintf("scale=%d:%d", w, h),
					"-c:a",
					"copy",
					"-c:v",
					"libx264",
					"-preset",
					"medium",
					"-crf",
					"23",
					"-y",
					dest)
				err = cmd.Run()
				if err != nil {
					global.Logger.Errorf("视频: %s :转码 %d*%d 失败。command : %s ,err info :%s", inputFile, w, h, cmd, err)
					continue
				}
				fmt.Println("转码后保存位置：", dest)
				src, _ := json.Marshal(common.Img{
					Src: dest,
					Tp:  "local",
				})
				switch r {
				case 1080:
					videoContribution.Video = src
				case 720:
					videoContribution.Video720p = src
				case 480:
					videoContribution.Video480p = src
				case 360:
					videoContribution.Video360p = src
				}
				if videoContribution.Save() {
					global.Logger.Errorf("视频 :%s : 转码%d*%d后视频保存到数据库失败", inputFile, w, h)
				}
				global.Logger.Infof("视频 :%s : 转码%d*%d成功", inputFile, w, h)
			}
		} else if data.VideoUploadType == "aliyunOss" && global.Config.AliyunOss.IsOpenTranscoding {
			// todo: 媒体服务已经开始收费没有免费体验额度，目前无法测试转码
			//wg.Add(1)
			// 3.2 oss视频使用媒资转码
			inputFile := data.VideoPath
			sr := strings.Split(inputFile, ".")
			// 云转码
			for _, r := range resolution {
				// 获取转码模板
				var template string
				dest := sr[0] + fmt.Sprintf("_output_%dp."+sr[1], r)
				src, _ := json.Marshal(common.Img{
					Src: dest,
					Tp:  data.VideoUploadType,
				})
				switch r {
				case 1080:
					template = global.Config.AliyunOss.TranscodingTemplate1080p
					videoContribution.Video = src
				case 720:
					template = global.Config.AliyunOss.TranscodingTemplate720p
					videoContribution.Video720p = src
				case 480:
					template = global.Config.AliyunOss.TranscodingTemplate480p
					videoContribution.Video480p = src
				case 360:
					template = global.Config.AliyunOss.TranscodingTemplate360p
					videoContribution.Video360p = src
				}
				outputUrl, _ := conversion.SwitchIngStorageFun(data.VideoUploadType, dest)
				taskName := "转码 : " + *data.Media + "时间 :" + time.Now().Format("2006.01.02 15:04:05") + " template : " + template
				jobInfo, err := oss.SubmitTranscodeJob(taskName, video.MediaID, outputUrl, template)
				if err != nil {
					global.Logger.Errorf("视频云转码 : %s 失败 err : %s", outputUrl, err.Error())
					continue
				}
				task := sundry.TranscodingTask{
					TaskID:     *jobInfo.TranscodeParentJob.ParentJobId,
					VideoID:    video.ID,
					Resolution: r,
					Dst:        dest,
					Status:     0,
					Type:       sundry.Aliyun,
				}
				if !task.AddTask() {
					global.Logger.Errorf("视频云转码任务名: %s 后将视频任务 保存到数据库失败", taskName)
				}
			}
		}
	}(width, height, videoContribution)
	wg.Wait()
	//fmt.Println("转码finish！")
	return videoContribution.ID, nil
}

// UpdateVideoContribution 更新视频信息
func UpdateVideoContribution(data *video.UpdateVideoContributionReceiveStruct, uid uint) (results interface{}, err error) {
	videoInfo := new(video3.VideosContribution)
	err = videoInfo.FindByID(data.ID)
	if err != nil {
		return nil, fmt.Errorf("更新视频不存在")
	}
	// 判读那这个视频是不是这个用户发布的
	if videoInfo.Uid != uid {
		return nil, fmt.Errorf("非法操作")
	}
	// 将封面img信息转为json串，存在数据库中，需要用的时候再Unmarshal为结构体
	coverImg, _ := json.Marshal(common.Img{
		Src: data.Cover,
		Tp:  data.CoverUploadType,
	})
	updateList := map[string]interface{}{
		"cover":     coverImg,
		"title":     data.Title,
		"label":     strings.Join(data.Label, ","),
		"reprinted": conversion.Bool2Int8(*data.Reprinted),
		"introduce": data.Introduce,
	}
	// 更新视频信息
	if !videoInfo.Update(updateList) {
		return nil, fmt.Errorf("更新视频数据失败")
	}
	return "视频信息更新成功", nil
}

// DeleteVideoByID 通过ID删除视频
// 1.从数据库查视频信息 2.在本地/oss删除视频 3.删除数据库记录
// todo:将一个视频对应的每个分辨率版本的视频其全部删除
// 验证的时候需要上传一个1080p的视频，去看oss中是否会有其他版本的视频，
// 然后执行这个删除操作，看OSS中是否全部删除
func DeleteVideoByID(data *video.DeleteVideoByIDReceiveStruct, uid uint) (results interface{}, err error) {
	v := new(video3.VideosContribution)

	err = v.FindByID(data.ID)
	if err != nil {
		global.Logger.Errorf("数据库查找视频信息失败：%v", err)
	}

	var deleteVideoPaths []string
	if v.Video720p != nil {
		//vo.Video720p.String()
		deleteVideoPaths = append(deleteVideoPaths, v.Video720p.String())
	} else if v.Video != nil {
		deleteVideoPaths = append(deleteVideoPaths, v.Video.String())
	} else if v.Video480p != nil {
		deleteVideoPaths = append(deleteVideoPaths, v.Video480p.String())
	} else if v.Video360p != nil {
		deleteVideoPaths = append(deleteVideoPaths, v.Video360p.String())
	} else {
		return nil, fmt.Errorf("视频不存在")
	}
	// todo:这里取出的v.Video720p是一个Json类型数据吧{Src:, Tp:}吧，怎么直接作为path放进数组了？ 打印一下看看
	// fmt.Println(deleteVideoPaths)

	if !v.Delete(data.ID, uid) {
		return nil, fmt.Errorf("删除视频失败")
	}
	// 删除OSS中的视频
	if err = oss.DeleteOSSFile(deleteVideoPaths); err != nil {
		return nil, fmt.Errorf("删除oss中视频失败：%v", err)
	}
	return "删除成功", nil
}

// VideoPostComment 发布视频评论
func VideoPostComment(data *video.VideosPostCommentReceiveStruct, uid uint) (results interface{}, err error) {
	node := filter.NewNode()
	node.StartMatch(data.Content)
	if node.IsSensitive() {
		return nil, fmt.Errorf("评论中包含敏感词")
	}
	videoInfo := new(video3.VideosContribution)
	err = videoInfo.FindByID(data.VideoID)
	if err != nil {
		return nil, errors.New("视频不存在")
	}
	// 找到被评论的那条评论
	ct := comments.Comment{}
	ct.ID = data.ContentID

	// 获取被评论那条评论的根评论
	CommentFirstID := ct.GetCommentFirstID()

	// 找到被评论那条评论的所属用户
	ctu := comments.Comment{}
	ctu.ID = data.ContentID
	commentUserID := ctu.GetCommentUserID()

	comment := comments.Comment{
		Uid:            uid,
		VideoID:        data.VideoID,
		Context:        data.Content,
		CommentID:      data.ContentID, // 父评论
		CommentUserID:  commentUserID,  // 父评论的用户
		CommentFirstID: CommentFirstID, // 根评论
	}
	if !comment.Create() {
		return nil, fmt.Errorf("发布视频评论失败")
	}

	// socket 推送(视频发布者在线情况下)
	if _, ok := users.Severe.UserMapChannel[videoInfo.UserInfo.ID]; ok {
		userChannel := users.Severe.UserMapChannel[videoInfo.UserInfo.ID]
		userChannel.NoticeMessage(notice.VideoComment)
	}
	return "发布视频评论成功", nil
}

// GetVideoManagementList 从创作中心获取视频稿件列表
func GetVideoManagementList(data *video.GetVideoManagementListReceiveStruct, uid uint) (results interface{}, err error) {
	// 获取个人发布视频信息
	list := new(video3.VideosContributionList)
	err = list.GetVideoManagementList(data.PageInfo, uid)
	if err != nil {
		return nil, fmt.Errorf("查询已发布视频列表失败")
	}
	res, err := video2.GetVideoManagementListResponse(list)
	if err != nil {
		return nil, fmt.Errorf("响应失败")
	}
	return res, nil
}

// LikeVideo 给视频点赞
func LikeVideo(data *video.LikeVideoReceiveStruct, uid uint) (results interface{}, err error) {
	videoInfo := new(video3.VideosContribution)
	err = videoInfo.FindByID(data.VideoID)
	if err != nil {
		return nil, fmt.Errorf("视频不存在")
	}
	lk := new(video3.Likes)
	err = lk.Like(uid, data.VideoID, videoInfo.UserInfo.ID)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("视频点赞操作失败")
	}

	// 视频作者在线时，向其推送信息
	if _, ok := users.Severe.UserMapChannel[videoInfo.UserInfo.ID]; ok {
		userChannel := users.Severe.UserMapChannel[videoInfo.UserInfo.ID]
		userChannel.NoticeMessage(notice.VideoLike)
	}
	return "视频点赞成功", nil
}

// DeleteVideoByPath 通过路径删除视频
func DeleteVideoByPath(data *video.DeleteVideoByPathReceiveStruct) (results interface{}, err error) {
	err = oss.DeleteOSSFile([]string{data.Path})
	if err != nil {
		return nil, fmt.Errorf("删除oss对象失败：%v", err)
	}
	return "删除视频成功", nil
}

// GetLastWatchTime 返回上次观看视频的进度
func GetLastWatchTime(uid, vid uint) (results interface{}, err error) {
	r := new(video3.WatchRecord)
	if err := r.GetByUidAndVideoId(uid, vid); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询视频观看进度失败：%v", err)
	}
	return r.WatchTime, nil
}

// SendWatchTime 保存视频观看进度
func SendWatchTime(data *video.SendWatchTimeReqStruct, uid uint) error {
	vid, _ := strconv.ParseUint(data.Id, 10, 64)
	r := new(video3.WatchRecord)
	global.MysqlDb.Model(&video3.WatchRecord{}).
		Where("uid = ? AND video_id = ?", uid, data.Id).
		First(r)
	if r.Id != 0 {
		r.WatchTime = data.Time
		return global.MysqlDb.Model(&video3.WatchRecord{}).
			Where("uid = ? AND video_id = ?", uid, data.Id).
			Updates(r).Error

	}
	r.Uid = uid
	r.VideoID = uint(vid)
	r.WatchTime = data.Time
	r.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	return global.MysqlDb.Model(&video3.WatchRecord{}).Save(&r).Error
}
