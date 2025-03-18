package config

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
	"path/filepath"
)

type KafkaConfigStruct struct {
	Server      string `ini:"server"`
	Brokers     string `ini:"brokers"`
	NormalTopic string `ini:"normalTopic"`
	DelayTopic  string `ini:"delayTopic"`
}

type SqlConfigStruct struct {
	IP       string `ini:"ip"`
	Port     int    `ini:"port"`
	Host     int    `ini:"host"`
	Username string `ini:"username"`
	Password string `ini:"password"`
	Database string `ini:"database"`
}

type RedisConfigStruct struct {
	IP       string `ini:"ip"`
	Port     int    `ini:"port"`
	Password string `ini:"password"`
}

type EmailConfigStruct struct {
	User     string `ini:"user"`
	Password string `ini:"password"`
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
}

type LiveConfigStruct struct {
	IP        string `ini:"ip"`
	Agreement string `ini:"agreement"`
	RTMP      string `ini:"rtmp"`
	FLV       string `ini:"flv"`
	HLS       string `ini:"hls"`
	Api       string `ini:"api"`
}

type ProjectConfigStruct struct {
	ProjectStates bool   `ini:"project_states"`
	Url           string `ini:"url"`
	UrlTest       string `ini:"url_test"`
}

type AliyunOss struct {
	Region                   string `ini:"region"`
	Bucket                   string `ini:"bucket"`
	AccessKeyId              string `ini:"accessKeyId"`
	AccessKeySecret          string `ini:"accessKeySecret"`
	Host                     string `ini:"host"`
	Endpoint                 string `ini:"endpoint"`
	RoleArn                  string `ini:"roleArn"`
	RoleSessionName          string `ini:"roleSessionName"`
	DurationSeconds          int    `ini:"durationSeconds"`
	IsOpenTranscoding        bool   `ini:"isOpenTranscoding"`
	TranscodingTemplate360p  string `ini:"transcodingTemplate360p"`
	TranscodingTemplate480p  string `ini:"transcodingTemplate480p"`
	TranscodingTemplate720p  string `ini:"transcodingTemplate720p"`
	TranscodingTemplate1080p string `ini:"transcodingTemplate1080p"`
	OssEndPoint              string `ini:"OssEndPoint"`
}

type Info struct {
	SqlConfig     *SqlConfigStruct
	RedisConfig   *RedisConfigStruct
	EmailConfig   *EmailConfigStruct
	LiveConfig    *LiveConfigStruct
	ProjectConfig *ProjectConfigStruct
	AliyunOss     *AliyunOss
	KafkaConfig   *KafkaConfigStruct
	ProjectUrl    string
}

var (
	Config = new(Info)
	cfg    *ini.File
	err    error
)

// getConfigPath 返回配置文件的路径
func getConfigPath() string {
	// 判断是测试还是正常运行，返回不同配置文件的路径
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if filepath.Base(curDir) == "test" {
		return filepath.ToSlash("../config/config.ini")
	}
	return filepath.ToSlash("./config/config.ini")
}

// 每个包的init函数都会在包被倒入的时候自动执行，即使该包被多个包导入也只会执行一次
func init() {
	path := getConfigPath()
	cfg, err = ini.Load(path)
	if err != nil {
		log.Fatalf("配置文件不存在，请检查环境：%v \n", err)
	}

	Config.SqlConfig = &SqlConfigStruct{}
	err = cfg.Section("mysql").MapTo(Config.SqlConfig)
	if err != nil {
		log.Fatalf("MySQL读取配置文件错误：%v \n", err)
	}

	Config.EmailConfig = &EmailConfigStruct{}
	err = cfg.Section("email").MapTo(Config.EmailConfig)
	if err != nil {
		log.Fatalf("Email读取配置文件错误：%v \n", err)
	}

	Config.KafkaConfig = &KafkaConfigStruct{}
	err = cfg.Section("kafka").MapTo(Config.KafkaConfig)
	if err != nil {
		log.Fatalf("Kafka读取配置文件错误：%v \n", err)
	}

	Config.RedisConfig = &RedisConfigStruct{}
	err = cfg.Section("redis").MapTo(Config.RedisConfig)
	if err != nil {
		log.Fatalf("Redis读取配置文件错误：%v \n", err)
	}

	Config.ProjectConfig = &ProjectConfigStruct{}
	err = cfg.Section("project").MapTo(Config.ProjectConfig)
	if err != nil {
		log.Fatalf("project读取配置文件错误：%v \n", err)
	}

	Config.LiveConfig = &LiveConfigStruct{}
	err = cfg.Section("live").MapTo(Config.LiveConfig)
	if err != nil {
		log.Fatalf("Live读取配置文件错误：%v \n", err)
	}

	Config.AliyunOss = &AliyunOss{}
	err = cfg.Section("aliyunOss").MapTo(Config.AliyunOss)
	if err != nil {
		log.Fatalf("AliyunOss读取配置文件错误：%v \n", err)
	}

	// 判断是否为测试环境，给定对应的url
	if Config.ProjectConfig.ProjectStates {
		Config.ProjectUrl = Config.ProjectConfig.Url
	} else {
		Config.ProjectUrl = Config.ProjectConfig.UrlTest
	}
}
