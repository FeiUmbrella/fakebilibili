package conversion

import (
	"encoding/json"
	"fakebilibili/infrastructure/model/common"
	"fakebilibili/infrastructure/pkg/global"
	"fmt"
	"strings"
)

// FormattingJsonSrc 返回图片的存储路径
func FormattingJsonSrc(str []byte) (url string, err error) {
	data := new(common.Img)
	err = json.Unmarshal(str, data)
	if err != nil {
		return "", fmt.Errorf("json format error")
	}
	if data.Src == "" {
		return "", nil
	}
	path, err := SwitchIngStorageFun(data.Tp, data.Src)
	if err != nil {
		return "", err
	}
	return path, nil
}

// SwitchIngStorageFun 根据类型拼接路径
func SwitchIngStorageFun(tp string, path string) (url string, err error) {
	prefix, err := SwitchTypeAsUrlPrefix(tp)
	//阿里云oss存储时，前缀是config里的host：https://easy-video-live.oss-cn-guangzhou.aliyuncs.com
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", prefix, path), nil
}

// SwitchTypeAsUrlPrefix 取url前缀
func SwitchTypeAsUrlPrefix(tp string) (url string, err error) {
	switch tp {
	case "local":
		return global.Config.ProjectUrl, nil
	case "aliyunOss":
		return global.Config.AliyunOss.Host, nil
	//todo:前端请求的coverUrl是在这里拼接的
	case "oss":
		//return "https://haifengonline-hangzhou.oss-cn-hangzhou.aliyuncs.com", nil
		return global.Config.AliyunOss.Host, nil
	case "wx":
		return "", nil
	default:
		return "", fmt.Errorf("undefined format")
	}
}

// StringConversionMap 将标签字符串拆分为单个标签
func StringConversionMap(s string) []string {
	s = strings.TrimSpace(s)
	list := strings.Split(s, ",")
	if len(list) == 1 && list[0] == "" {
		return []string{}
	}
	return list
}
