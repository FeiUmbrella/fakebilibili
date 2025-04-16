package osscommonality

import "fakebilibili/infrastructure/model/common"

// UploadSliceInfo 文件分片的下标以及分片内容hash
type UploadSliceInfo struct {
	Index int    `json:"index" `
	Hash  string `json:"hash"`
}
type UploadSliceList []UploadSliceInfo

// UploadCheckStruct 文件内容hash作为文件的name
type UploadCheckStruct struct {
	FileMd5   string          `json:"file_md5" binding:"required"`
	Interface string          `json:"interface" binding:"required"`
	SliceList UploadSliceList `json:"slice_list"  binding:"required"` // 文件分成的分片列表
}

// UploadMergeStruct 将文件的所有分片合并为原文件
type UploadMergeStruct struct {
	FileName  string          `json:"file_name" binding:"required"`   // 原文件名
	Interface string          `json:"interface" binding:"required"`   // 保存接口/也可以所示文件的类型
	SliceList UploadSliceList `json:"slice_list"  binding:"required"` // 文件所有分片的下标和hash name
}

type UploadingMethodStruct struct {
	Method string `json:"method"  binding:"required"`
}

type UploadingDirStruct struct {
	Interface string `json:"interface"  binding:"required"`
}

type GetFullPathOfImageMethodStruct struct {
	Path string `json:"path"  binding:"required"`
	Type string `json:"type"  binding:"required"`
}

type SearchStruct struct {
	PageInfo common.PageInfo `json:"page_info" binding:"required"`
	Type     string          `json:"type" binding:"required"`
}

type RegisterMediaStruct struct {
	Type string `json:"type"` //aliyunOss
	Path string `json:"path"` // E:/video/hash.mp4
}
