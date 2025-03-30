package users

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fakebilibili/adapter/http/receive"
	"fakebilibili/adapter/http/response"
	"fakebilibili/infrastructure/consts"
	"fakebilibili/infrastructure/model/common"
	user2 "fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/model/user/attention"
	"fakebilibili/infrastructure/model/user/chat"
	"fakebilibili/infrastructure/model/user/collect"
	"fakebilibili/infrastructure/model/user/favorites"
	"fakebilibili/infrastructure/model/user/liveInfo"
	"fakebilibili/infrastructure/model/user/notice"
	"fakebilibili/infrastructure/model/user/record"
	"fakebilibili/infrastructure/pkg/global"
	"fakebilibili/infrastructure/pkg/utils/conversion"
	"fakebilibili/infrastructure/pkg/utils/email"
	"fakebilibili/infrastructure/pkg/utils/jwt"
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// GetUserInfo 获取用户信息
func GetUserInfo(uid uint) (interface{}, error) {
	user := new(user2.User)
	user.IsExistByField("id", uid)
	res := response.GetUserInfoResponse(user)
	return res, nil
}

// SetUserInfo 设置用户信息
func SetUserInfo(data *receive.SetUserInfoReceiveStruct, uid uint) (interface{}, error) {
	fields := map[string]interface{}{
		"Username":    data.Username,
		"Gender":      data.Gender,
		"BirthDate":   data.BirthDate,
		"IsVisible":   conversion.Bool2Int8(*data.IsVisible),
		"Signature":   data.Signature,
		"SocialMedia": data.SocialMedia,
	}
	user := new(user2.User)
	user.ID = uid
	return user.UpdatePureZero(fields), nil
}

// DetermineNameExists 判断名字是否存在
func DetermineNameExists(data *receive.DetermineNameExistsStruct, uid uint) (interface{}, error) {
	user := new(user2.User)
	exist := user.IsExistByField("username", data.Username)
	if user.ID == uid { // 新名字和原来一样
		return false, nil
	} else if exist {
		return true, nil
	}
	return false, nil
}

// UpdateAvatar 更新头像
func UpdateAvatar(data *receive.UpdateAvatarStruct, uid uint) (interface{}, error) {
	user := new(user2.User)
	user.ID = uid
	photo, _ := json.Marshal(common.Img{
		Src: data.ImgUrl,
		Tp:  data.Tp,
	})
	user.Photo = photo
	if user.Update() {
		return conversion.SwitchIngStorageFun(data.Tp, data.ImgUrl)
	}
	return nil, fmt.Errorf("更新头像失败")
}

// GetLiveInfo 获取直播间信息
func GetLiveInfo(uid uint) (interface{}, error) {
	info := new(liveInfo.LiveInfo)
	if info.IsExistByField("uid", uid) {
		res, err := response.GetLiveInfoResponse(info)
		if err != nil {
			return nil, fmt.Errorf("获取直播间信息失败")
		}
		return res, nil
	}
	return common.Img{}, nil
}

// SaveLiveInfo 修改直播间信息
func SaveLiveInfo(data *receive.SaveLiveDataReceiveStruct, uid uint) (interface{}, error) {
	img, _ := json.Marshal(common.Img{
		Src: data.ImgUrl,
		Tp:  data.Tp,
	})
	info := &liveInfo.LiveInfo{
		Uid:   uid,
		Title: data.Title,
		Img:   datatypes.JSON(img),
	}
	if info.UpdateInfo() {
		return "直播间信息修改成功", nil
	} else {
		return nil, fmt.Errorf("修改失败")
	}
}

// SendEmailVerificationCodeByChangePassword 在登陆状态下修改密码
func SendEmailVerificationCodeByChangePassword(uid uint) (interface{}, error) {
	user := new(user2.User)
	user.Find(uid)

	// 发送对象
	mailTo := []string{user.Email}
	// 邮箱主题
	subject := "验证码"
	// 生成6位验证码
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(1000000))
	// 邮件正文
	body := fmt.Sprintf("您正在修改密码，您的验证码为:%s,5分钟内有效,请勿转发他人", code)
	err := email.SendEmail(mailTo, subject, body)
	if err != nil {
		return nil, err
	}
	err = global.RedisDb.Set(fmt.Sprintf("%s%s", consts.EmailVerificationCodeByChangePassword, user.Email), code, 5*time.Minute).Err()
	if err != nil {
		return nil, err
	}
	return "正在修改密码，邮箱验证码发送成功", nil
}

// ChangePassword 登录状态下修改密码
func ChangePassword(data *receive.ChangePasswordReceiveStruct, uid uint) (interface{}, error) {
	user := new(user2.User)
	user.Find(uid)

	if data.Password != data.ConfirmPassword {
		return nil, fmt.Errorf("两次密码不一致！")
	}

	// 判断邮箱验证码是否一致
	verCode, err := global.RedisDb.Get(fmt.Sprintf("%s%s", consts.EmailVerificationCodeByChangePassword, user.Email)).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("验证码过期")
	}
	if verCode != data.VerificationCode {
		return nil, fmt.Errorf("验证码错误")
	}

	// 更新密码
	// 生成6位密码盐
	salt := make([]byte, 6)
	for i := range salt {
		salt[i] = jwt.SaltStr[rand.Int63()%int64(len(jwt.SaltStr))]
	}
	// p = salt + password + salt
	password := []byte(fmt.Sprintf("%s%s%s", salt, data.Password, salt))
	// md5 加密
	passwordMd5 := fmt.Sprintf("%x", md5.Sum(password))

	user.Salt = string(salt)
	user.Password = passwordMd5

	res := user.Update()
	if !res {
		return nil, fmt.Errorf("密码修改失败")
	}
	return "密码修改成功", nil
}

// Attention 关注/取消关注
func Attention(data *receive.AttentionReceiveStruct, uid uint) (interface{}, error) {
	at := new(attention.Attention)
	if uid == data.Uid {
		return nil, fmt.Errorf("用户uid:%x 关注 用户uid:%x, 失败", uid, data.Uid)
	}
	if at.Attention(uid, data.Uid) {
		return "关注成功", nil
	}
	return nil, fmt.Errorf("关注失败")
}

// CreateFavorites 创建或更新收藏夹
func CreateFavorites(data *receive.CreateFavoritesReceiveStruct, uid uint) (interface{}, error) {
	// 当传进的收藏夹ID<=0的时候，表明是要创建一个新的收藏夹；否则表明是更新某个收藏夹
	if data.ID <= 0 { // 创建一个新的收藏夹
		if len(data.Title) == 0 {
			return nil, fmt.Errorf("收藏夹标题为空")
		}
		// 只有标题
		if len(data.Tp) == 0 && len(data.Content) == 0 && len(data.Cover) == 0 {
			fs := &favorites.Favorites{
				Uid:   uid,
				Title: data.Title,
				Max:   1000,
			}
			if !fs.Create() {
				return nil, fmt.Errorf("收藏夹创建失败")
			}
		} else {
			cover, _ := json.Marshal(common.Img{
				Src: data.Cover,
				Tp:  data.Tp,
			})
			fs := &favorites.Favorites{
				Uid:     uid,
				Title:   data.Title,
				Cover:   cover,
				Content: data.Content,
				Max:     1000,
			}
			if !fs.Create() {
				return nil, fmt.Errorf("收藏夹创建失败")
			}
		}
		return "更新创建成功", nil
	} else {
		// 对已有的收藏夹进行更新
		fs := &favorites.Favorites{}
		if !fs.Find(data.ID) {
			return nil, fmt.Errorf("查询失败")
		}
		if fs.Uid != uid { // 查询的收藏夹不是本用户的收藏夹
			return nil, fmt.Errorf("查询非法操作")
		}
		cover, _ := json.Marshal(common.Img{
			Src: data.Cover,
			Tp:  data.Tp,
		})
		fs.Cover = cover
		fs.Title = data.Title
		fs.Content = data.Content
		if !fs.Update() {
			return nil, fmt.Errorf("更新收藏夹失败")
		}
		return "更新创建成功", nil
	}
}

// GetFavoritesList 以列表形式获取收藏夹
func GetFavoritesList(uid uint) (interface{}, error) {
	fs := new(favorites.FavoriteList)
	err := fs.GetFavoritesList(uid)
	if err != nil {
		return nil, fmt.Errorf("获取所有收藏夹失败")
	}
	res, err := response.GetFavoritesListResponse(fs)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteFavorites 删除收藏夹
func DeleteFavorites(data *receive.DeleteFavoritesReceiveStruct, uid uint) (interface{}, error) {
	fs := new(favorites.Favorites)
	err := fs.Delete(data.ID, uid)
	if err != nil {
		return nil, err
	}
	return "收藏夹删除成功", nil
}

// FavoriteVideo 收藏/取消收藏视频
// // 设VideoID原来所在的收藏夹的id集合为old_ids，如果IDs中id不在old_ids中，那么是要收藏该视频到为id的收藏夹中
// // 如果old_ids中id不在IDs中，那么是要取消收藏在id的收藏夹中的该视频
func FavoriteVideo(data *receive.FavoriteVideoReceiveStruct, uid uint) (interface{}, error) {
	// 查看是否需要取消某个收藏夹的收藏
	cl := new(collect.CollectsList)
	err := cl.FindFavoriteIncludeVideo(data.VideoID) // 查找目前收藏该视频的所有收藏夹
	if err != nil {
		return nil, fmt.Errorf("取消视频收藏时，查询所在收藏夹失败")
	}
	oldIds := make([]uint, 0)
	for _, v := range *cl {
		if v.Uid != uid {
			continue
		} // 不是该用户的收藏夹跳过记录
		oldIds = append(oldIds, v.FavoritesID)
	}
	for _, id := range oldIds {
		// 判断在新的收藏夹id中是否包含当前id，不包含则需要取消收藏
		if !idInIds(id, data.IDs) {
			deleteCollect := new(collect.Collect)
			err := deleteCollect.DeleteOneVideoInFavorite(data.VideoID, id)
			if err != nil {
				return nil, fmt.Errorf("取消视频收藏时，取消收藏失败")
			}
		}
	}

	// 查看哪些收藏夹需要收藏该视频
	for _, id := range data.IDs {
		fs := new(favorites.Favorites)
		fs.Find(id)
		fmt.Printf("用户id：%x,收藏夹：%x,收藏夹所属用户：%x", uid, fs.ID, fs.Uid)
		// 这个收藏夹不属于该用户
		if fs.Uid != uid {
			return nil, fmt.Errorf("非法操作，没有该收藏夹权限")
		}
		if len(fs.CollectList)+1 > fs.Max {
			return nil, fmt.Errorf("收藏夹已满")
		}

		// 查看是否重复收藏
		videoCollect := new(collect.Collect)
		err = videoCollect.FindOneVideoInFavorite(data.VideoID, fs.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) { // 没有找到，收藏该视频
				videoCreate := &collect.Collect{
					Uid:         uid,
					FavoritesID: id,
					VideoID:     data.VideoID,
				}
				if !videoCreate.Create() {
					return nil, fmt.Errorf("视频收藏失败")
				}
			}
			return nil, fmt.Errorf("收藏视频时，视频查重失败")
		}
		// 这里表明找到已添加的记录，不作操作
	}
	return "视频收藏/取消收藏操作完成", nil
}

// idInIds 判断id是否在Ids
func idInIds(id uint, ids []uint) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}

// GetFavoritesListByFavoriteVideo 获取用户包含某个视频的收藏夹列表
func GetFavoritesListByFavoriteVideo(data *receive.GetFavoritesListByFavoriteVideoReceiveStruct, uid uint) (interface{}, error) {
	// 先找到用户所有收藏夹
	fl := new(favorites.FavoriteList)
	err := fl.GetFavoritesList(uid)
	if err != nil {
		return nil, fmt.Errorf("获取用户所有收藏夹失败")
	}

	// 找到包含目标视频的所有收藏夹fid
	cl := new(collect.CollectsList)
	err = cl.FindFavoriteIncludeVideo(data.VideoID)
	if err != nil {
		return nil, fmt.Errorf("获取包含目标视频所有收藏夹失败")
	}
	fid := make([]uint, 0)
	for _, v := range *cl {
		fid = append(fid, v.FavoritesID)
	}

	// 处理返回
	res, err := response.GetFavoritesListByFavoriteVideoResponse(fl, fid)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetFavoriteVideoList 获取收藏夹中的视频列表
func GetFavoriteVideoList(data *receive.GetFavoriteVideoListReceiveStruct) (interface{}, error) {
	cl := new(collect.CollectsList)
	err := cl.FindVideosByFavoriteID(data.FavoriteID)
	if err != nil {
		return nil, fmt.Errorf("查询收藏夹视频列表失败")
	}

	res, err := response.GetFavoriteVideoListResponse(cl)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetCollectListName 根据收藏夹id获取收藏夹的title
func GetCollectListName(data *receive.GetCollectListNameReceiveStruct) (interface{}, error) {
	f := new(favorites.Favorites)
	f.FindFavoriteByFID(data.FavoriteID)
	return *f, nil
}

// GetRecordList 获取历史记录
func GetRecordList(data *receive.GetRecordListReceiveStruct, uid uint) (interface{}, error) {
	rL := new(record.RecordList)
	err := rL.GetRecordListByUid(uid, data.PageInfo)
	if err != nil {
		return nil, fmt.Errorf("查询历史记录失败")
	}
	res, err := response.GetRecordListResponse(rL)
	if err != nil {
		return nil, fmt.Errorf("返回历史记录失败")
	}
	return res, nil
}

// ClearRecord 清空历史记录
func ClearRecord(uid uint) (interface{}, error) {
	rd := new(record.Record)
	err := rd.ClearRecord(uid)
	if err != nil {
		return nil, fmt.Errorf("清空记录失败")
	}
	return "清空历史记录成功", nil
}

// DeleteRecordByID 删除一条历史记录
func DeleteRecordByID(data *receive.DeleteRecordByIDReceiveStruct, uid uint) (interface{}, error) {
	rd := new(record.Record)
	err := rd.DeleteRecordByID(data.ID, uid)
	if err != nil {
		return nil, fmt.Errorf("删除历史记录失败")
	}
	return "删除一条历史记录成功", nil
}

// GetNoticeList 获取通知列表
func GetNoticeList(data *receive.GetNoticeListReceiveStruct, uid uint) (interface{}, error) {
	msgType := make([]string, 0)
	noticeList := new(notice.NoticesList)
	switch data.Type {
	case "comment":
		msgType = append(msgType, notice.VideoComment, notice.ArticleComment)
		break
	case "like":
		msgType = append(msgType, notice.VideoLike, notice.ArticleLike)
		break
	case "system":
		// todo:这里为什么插入两个DailyReport
		msgType = append(msgType, notice.UserLogin, notice.DailyReport, notice.DailyReport)
		break
	}

	err := noticeList.GetNoticeList(data.PageInfo, msgType, uid)
	if err != nil {
		return nil, fmt.Errorf("查询失败")
	}
	// 全部通知设为已读
	nt := new(notice.Notice)
	err = nt.ReadAll(uid)
	if err != nil {
		return nil, fmt.Errorf("设置通知消息为已读失败")
	}
	res, err := response.GetNoticeListResponse(noticeList)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetChatList 获取消息列表
func GetChatList(uid uint) (interface{}, error) {
	// 获取消息列表
	chatList := new(chat.ChatList)
	err := chatList.GetListByID(uid)
	if err != nil {
		return nil, fmt.Errorf("查询用户消息列表失败")
	}
	toIDs := make([]uint, 0) // 发送对象
	for _, c := range *chatList {
		toIDs = append(toIDs, c.Tid)
	}
	msgList := make(map[uint]*chat.MsgList)
	for _, toID := range toIDs {
		mList := new(chat.MsgList)
		err = mList.FindList(uid, toID) // 查找两人的聊天信息
		if err != nil {
			break
		}
		msgList[toID] = mList
	}
	res, err := response.GetChatListResponse(chatList, msgList)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetChatHistoryMsg 获取历史聊天记录
func GetChatHistoryMsg(data *receive.GetChatHistoryMsgStruct, uid uint) (interface{}, error) {
	ml := new(chat.MsgList)
	err := ml.FindHistoryMst(uid, data.Tid, data.LastTime)
	if err != nil {
		return nil, fmt.Errorf("查询历史聊天记录失败")
	}
	res, err := response.GetChatHistoryMsgResponse(ml)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// PersonalLetter 点击“私信”时触发。
// 如果两人曾聊过天那么加载最后一条信息，否则在聊天列表新建两人聊天（因为只是点击“私信”而非发送信息，此时聊天记录为空）
func PersonalLetter(data *receive.PersonalLetterReceiveStruct, uid uint) (interface{}, error) {
	msg := new(chat.Msg)
	err := msg.GetLastMessage(uid, data.ID)
	// gorm.First的error返回值：1.gorm.ErrRecordNotFound没有对应记录; 2. nil找到对应记录;3.其他错误
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查找私信最后一条信息操作失败")
	}
	var lastTime time.Time
	if err == nil { // 找到一条记录
		lastTime = msg.CreatedAt
	} else { // 两人没有聊过天
		lastTime = time.Now()
	}
	// 在私信列表创建一个聊天框
	cl := &chat.ChatsListInfo{
		Uid:         uid,
		Tid:         data.ID,
		LastMessage: msg.Message, // 如果从没聊过天那么msg在数据库也没找到对应记录，这个是对应零值
		LastAt:      lastTime,
	}
	err = cl.AddChat()
	if err != nil {
		return nil, fmt.Errorf("创建私信列表失败")
	}
	return "操作成功", nil
}

// DeleteChatItem 删除聊天记录
func DeleteChatItem(data *receive.DeleteChatItemReceiveStruct, uid uint) (interface{}, error) {
	chatInfo := new(chat.ChatsListInfo)
	err := chatInfo.DeleteChat(data.ID, uid)
	if err != nil {
		return nil, fmt.Errorf("删除聊天记录失败")
	}
	return "删除聊天记录成功", nil
}
