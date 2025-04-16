package osscommonality

import (
	osscommonality2 "fakebilibili/adapter/http/receive/osscommonality"
	"fakebilibili/adapter/http/response/osscommonality"
	"fakebilibili/infrastructure/model/contribution/video"
	"fakebilibili/infrastructure/model/sundry"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/model/user/attention"
	"fakebilibili/infrastructure/pkg/global"
	"fakebilibili/infrastructure/pkg/utils/conversion"
	"fakebilibili/infrastructure/pkg/utils/location"
	"fakebilibili/infrastructure/pkg/utils/oss"
	"fakebilibili/infrastructure/pkg/utils/validator"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	// Temporary 文件切片保存目录
	Temporary = filepath.ToSlash("assets/temp")
)

// OssSTS 获取STS临时授权码，这个方法返回500的原因是config.ini配置有问题
func OssSTS() (interface{}, error) {
	info, err := oss.GetStsInfo()
	if err != nil {
		global.Logger.Errorf("获取OssSts密钥失败 错误原因 :%s", err.Error())
		return nil, fmt.Errorf("获取失败")
	}
	res, err := osscommonality.GteStsInfo(info)
	if err != nil {
		return nil, fmt.Errorf("响应失败")
	}
	return res, nil
}

// Upload 保存上传的文件到本地
func Upload(file *multipart.FileHeader, ctx *gin.Context) (interface{}, error) {
	// 设置文件名、存储目录等基础信息
	// 如果文件大小超过maxMem bytes，则maxMem存储在内存，超过部分存储在磁盘的临时文件
	err := ctx.Request.ParseMultipartForm(128)
	if err != nil {
		return nil, err
	}
	mForm := ctx.Request.MultipartForm
	// 文件名
	// todo: 这里fileName为空，join的话会使得传进来的文件名直接相连吧？
	var fileName string
	fileName = strings.Join(mForm.Value["name"], fileName)
	// todo:这里fileInterface也同样为空，直接使用了Join？
	var fileInterface string
	fileInterface = strings.Join(mForm.Value["interface"], fileInterface)

	method := new(sundry.Upload)
	if !method.IsExistByField("interface", fileInterface) {
		return nil, fmt.Errorf("上传接口不存在")
	}
	if len(method.Path) == 0 {
		return nil, fmt.Errorf("请联系管理员设置接口保存路径")
	}
	// 检查文件后缀
	suffix := fileName[strings.LastIndex(fileName, "."):]
	err = validator.CheckVideoSuffix(suffix)
	if err != nil {
		return nil, fmt.Errorf("非法后缀")
	}
	if !location.IsDir(method.Path) {
		// 创建多级目录
		if err = os.MkdirAll(method.Path, 0775); err != nil {
			global.Logger.Errorf("创建文件报错路径失败 创建路径为：%s 错误原因 : %s", method.Path, err.Error())
			return nil, fmt.Errorf("创建保存路径失败")
		}
	}

	// 2.保存文件
	dest := filepath.ToSlash(method.Path + "/" + fileName)
	err = ctx.SaveUploadedFile(file, dest)
	if err != nil {
		global.Logger.Errorf("保存文件失败-保存路径为：%s ,错误原因 : %s", dest, err.Error())
		return nil, fmt.Errorf("上传失败")
	}
	return dest, nil
}

// UploadSlice 分片上传，将文件的一个分片保存在临时目录，以便后面进行合并操作
func UploadSlice(file *multipart.FileHeader, ctx *gin.Context) (interface{}, error) {
	err := ctx.Request.ParseMultipartForm(128)
	if err != nil {
		return nil, err
	}
	mForm := ctx.Request.MultipartForm
	// 文件名
	var fileName string
	fileName = strings.Join(mForm.Value["name"], fileName)
	var fileInterface string
	fileInterface = strings.Join(mForm.Value["interface"], fileInterface)

	method := new(sundry.Upload)
	if !method.IsExistByField("interface", fileInterface) {
		return nil, fmt.Errorf("上传接口不存在")
	}
	if len(method.Path) == 0 {
		return nil, fmt.Errorf("请联系管理员设置接口保存路径")
	}

	if !location.IsDir(Temporary) {
		if err = os.MkdirAll(method.Path, 0775); err != nil {
			global.Logger.Errorf("创建文件报错路径失败 创建路径为：%s 错误原因 : %s", Temporary, err.Error())
			return nil, fmt.Errorf("创建保存路径失败")
		}
	}

	// 2.保存文件
	dest := filepath.ToSlash(Temporary + "/" + fileName)
	err = ctx.SaveUploadedFile(file, dest)
	if err != nil {
		global.Logger.Errorf("文件分片失败-保存路径为：%s ,错误原因 : %s", dest, err.Error())
		return nil, fmt.Errorf("上传失败")
	}
	_ = os.Chmod(dest, 0775)
	return dest, nil
}

// UploadCheck 前端传来的分片标号列表，返回未保存到本地的分片标号
func UploadCheck(data *osscommonality2.UploadCheckStruct) (interface{}, error) {
	method := new(sundry.Upload)
	if !method.IsExistByField("interface", data.Interface) {
		return nil, fmt.Errorf("未配置上传方法")
	}

	list := make(osscommonality2.UploadSliceList, 0)
	path := filepath.ToSlash(method.Path + "/" + data.FileMd5) // 原文件的保存路径，不是切片
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		// 原文件已存在不用管分片处理了
		global.Logger.Infof("上传文件 %s 已存在", data.FileMd5)
		return osscommonality.UploadCheckResponse(true, list, path)
	}
	// 取出未上传的分片
	for _, slice := range data.SliceList {
		if _, err := os.Stat(filepath.ToSlash(Temporary + "/" + slice.Hash)); os.IsNotExist(err) {
			list = append(list, osscommonality2.UploadSliceInfo{
				Index: slice.Index,
				Hash:  slice.Hash,
			})
		}
	}
	return osscommonality.UploadCheckResponse(false, list, "")
}

// UploadMerge 将所有分片上传完毕的文件分片进行合并为原文件
func UploadMerge(data *osscommonality2.UploadMergeStruct) (interface{}, error) {
	method := new(sundry.Upload)
	if !method.IsExistByField("interface", data.Interface) {
		return nil, fmt.Errorf("未配置上传方法")
	}
	// 1.创建文件保存目录
	if !location.IsDir(filepath.ToSlash(method.Path)) {
		if err := os.MkdirAll(filepath.ToSlash(method.Path), 0775); err != nil {
			global.Logger.Errorf("创建文件报错路径失败 创建路径为：%s", method.Path)
			return nil, fmt.Errorf("创建保存路径失败")
		}
	}
	// 2. 文件保存路径
	dest := filepath.ToSlash(method.Path + "/" + data.FileName)
	list := make(osscommonality2.UploadSliceList, 0)
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		// 文件已存在直接返回
		return dest, nil
	}
	// 3.取出未上传分片
	for _, slice := range data.SliceList {
		if _, err := os.Stat(filepath.ToSlash(Temporary + "/" + slice.Hash)); os.IsNotExist(err) {
			list = append(list, osscommonality2.UploadSliceInfo{
				Index: slice.Index,
				Hash:  slice.Hash,
			})
		}
	}
	// 4.若有未上传分片，直接返回，让其分片上传完成后再进行合并操作
	if len(list) > 0 {
		global.Logger.Warnf("上传文件 %s 分片未全部上传", data.FileName)
		return nil, fmt.Errorf("分片未全部上传")
	}
	// 创建空原文件
	destFile, err := os.Create(dest)
	if err != nil {
		global.Logger.Errorf("创建的合并后文件失败 err : %s", err.Error())
	}
	if err = destFile.Close(); err != nil {
		global.Logger.Errorf("创建的合并后文件关闭失败 %d", err)
	}
	destFile, err = os.OpenFile(dest, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModeSetuid) // os.ModeSetuid 表示文件具有其创建者用户id权限
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			global.Logger.Errorf("关闭资源 err : %s", err)
		}
	}(destFile)

	// 合并分片到destFile
	for _, slice := range data.SliceList {
		sliceFile, err := os.OpenFile(filepath.ToSlash(Temporary+"/"+slice.Hash), os.O_RDONLY, os.ModePerm) // 0777
		if err != nil {
			global.Logger.Errorf("合并操作打开临时分片失败 错误原因 : %s", err)
			break
		}
		temp, err := io.ReadAll(sliceFile)
		if err != nil {
			global.Logger.Errorf("合并操作读取分片失败 错误原因 : %s", err)
			break
		}
		if _, err := destFile.Write(temp); err != nil {
			global.Logger.Errorf("合并分片追加错误 错误原因 : %s", err)
			return nil, fmt.Errorf("合并分片追加错误")
		}
		// 关闭分片
		_ = sliceFile.Close()
		if err = os.Remove(sliceFile.Name()); err != nil {
			global.Logger.Errorf("合并操作删除临时分片失败 错误原因 : %s", err)
		}
	}
	return dest, nil
}

func UploadingMethod(data *osscommonality2.UploadingMethodStruct) (results interface{}, err error) {
	method := new(sundry.Upload)
	// 从前端传进来的data.Method就是传的interface字段
	if method.IsExistByField("interface", data.Method) {
		return osscommonality.UploadingMethodResponse(method.Method), nil
	} else {
		return nil, fmt.Errorf("未配置上传方法")
	}
}

// UploadingDir 返回文件存储目录，已经图片文件的大小比如5MB
func UploadingDir(data *osscommonality2.UploadingDirStruct) (results interface{}, err error) {
	method := new(sundry.Upload)
	if method.IsExistByField("interface", data.Interface) {
		return osscommonality.UploadingDirResponse(method.Path, method.Quality), nil
	} else {
		return nil, fmt.Errorf("未配置上传方法")
	}
}

// GetFullPathOfImage 获取图片文件的保存路径
func GetFullPathOfImage(data *osscommonality2.GetFullPathOfImageMethodStruct) (results interface{}, err error) {
	path, err := conversion.SwitchIngStorageFun(data.Type, data.Path)
	if err != nil {
		return nil, err
	}
	return path, nil
}

// Search 根据关键词搜索用户/视频
func Search(data *osscommonality2.SearchStruct, uid uint) (interface{}, error) {
	switch data.Type {
	case "video":
		// 视频搜索
		list := new(video.VideosContributionList)
		err := list.Search(data.PageInfo)
		if err != nil {
			return nil, fmt.Errorf("查询失败")
		}
		res, err := osscommonality.SearchVideoResponse(list)
		if err != nil {
			return nil, fmt.Errorf("查询失败")
		}
		return res, nil
	case "user":
		list := new(user.UserList)
		err := list.Search(data.PageInfo) // 返回名字中包含keyword的用户
		if err != nil {
			return nil, fmt.Errorf("查询失败")
		}
		aids := make([]uint, 0)
		if uid != 0 {
			// 用户登入的话
			al := new(attention.AttentionsList)
			err = al.GetAttentionList(uid) // 获取关注列表
			if err != nil {
				global.Logger.Errorf("用户id %d 获取取关注列表失败,错误原因 : %s ", uid, err.Error())
				return nil, fmt.Errorf("获取关注列表失败")
			}
			for _, v := range *al {
				aids = append(aids, v.AttentionID)
			}
		}
		res, err := osscommonality.SearchUserResponse(list, aids)
		return res, nil
	default:
		return nil, fmt.Errorf("未匹配的类型")
	}
}

// RegisterMedia 将某个保存在OSS上的视频注册媒体资源，后续可以利用阿里云对该视频的一些功能
func RegisterMedia(data *osscommonality2.RegisterMediaStruct) (interface{}, error) {
	path, _ := conversion.SwitchIngStorageFun(data.Type, data.Path)
	// 注册媒资
	registerMediaBody, err := oss.RegisterMediaInfo(path, "video", time.Now().String())
	if err != nil {
		return nil, fmt.Errorf("注册媒资失败")
	}
	if registerMediaBody == nil {
		global.Logger.Infoln("注册媒体资源的返回结果为空")
		return nil, fmt.Errorf("注册媒体资源的返回结果为空")
	}
	return registerMediaBody.MediaId, nil
}
