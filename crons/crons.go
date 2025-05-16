package crons

import (
	"bufio"
	"encoding/json"
	"fakebilibili/domain/servicce/contribution/article"
	"fakebilibili/domain/servicce/contribution/video"
	"fakebilibili/infrastructure/model/cron_events"
	"fakebilibili/infrastructure/model/home"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/model/user/attention"
	"fakebilibili/infrastructure/model/user/notice"
	"fakebilibili/infrastructure/pkg/global"
	"fmt"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var job *cron.Cron

func InitCron() {
	/*
		这里可以启动消费者
	*/
	// 覆盖Cron默认解析格式，添加秒级
	job = cron.New(cron.WithSeconds())
	// 每天零点持久化runtime/log文件、发送日报、更新轮播图
	job.AddFunc("@midnight", StoreRuntimeLogFile)
	job.AddFunc("@midnight", DailyReport)
	job.AddFunc("@midnight", UpdateRotograph)

	job.Start()
}

// 备份前一天日志并创建下一天日志
var (
	files = []string{
		fmt.Sprintf("./runtime/log/%s/error.log", time.Now().Add(-1*time.Hour).Format(time.DateOnly)),
		fmt.Sprintf("./runtime/log/%s/info.log", time.Now().Add(-1*time.Hour).Format(time.DateOnly)),
	}
)

// StoreRuntimeLogFile 持久化日志文件
func StoreRuntimeLogFile() {
	for _, filePath := range files {
		if err := ProcessFile(filePath); err != nil {
			global.Logger.Errorf("持久化日志出错，文件路径为%s", filePath)
		}
	}
}

// ProcessFile 处理文件
func ProcessFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("File %s does not exist, skipping...", filePath)
			return nil
		}
		log.Fatal("打开文件出错：", err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		entry := new(cron_events.RuntimeLogEntry)
		if err := json.Unmarshal([]byte(line), entry); err != nil {
			log.Print("读取日志行出错" + err.Error())
			continue
		}
		if err := entry.Create(); err != nil { // 新建一条数据库日志记录
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}
	return nil
}

// DailyReport 用户日报
func DailyReport() {
	fmt.Println("开始发送日报！")
	var content string = "您昨日报告"
	users := new(user.UserList)
	err := users.GetAllUserIds()
	if err != nil {
		global.Logger.Errorln("获取所有用户err：" + err.Error())
	}

	// 昨日涨粉情况
	at := new(attention.Attention)
	for _, usr := range *users {
		count, err := at.GetNewAddAttentionByTime(time.Now().Add(-24*time.Hour).Format(time.DateTime), usr.ID)
		if err != nil {
			global.Logger.Errorln("获取新增粉丝数量err：" + err.Error())
			return
		}
		content = fmt.Sprintf("%s:昨日新增%d粉丝！", content, count)
		ne := new(notice.Notice)
		if err = ne.AddNotice(usr.ID, 2, 0, notice.DailyReport, content); err != nil {
			global.Logger.Errorln("发送日报出错err：" + err.Error())
			return
		}
	}
	global.Logger.Infoln("向所有用户发送日报成功！")
}

// UpdateRotograph 定时更新热门轮播图
func UpdateRotograph() {
	fmt.Println("开始更新轮播图！")
	var err error
	// 找到最热门的两个视频和一个专栏，放进rotograph表中
	heatestVideos, err := video.GetTop2HeatVideos()
	if err != nil {
		global.Logger.Errorln("查询最热2视频err：" + err.Error())
		return
	}
	heatestArticle, err := article.GetHeatestArticle()
	if err != nil {
		global.Logger.Errorln("查询最热1专栏err：" + err.Error())
		return
	}

	// 将查询结构插入Rotograph表
	var rotographs []*home.Rotograph
	for _, vd := range heatestVideos {
		rotographs = append(rotographs, &home.Rotograph{
			Title: vd.Title,
			Cover: vd.Cover,
			Color: "rgb(116,82,81)",
			Type:  "video",
			ToId:  vd.ID,
		})
	}
	rotographs = append(rotographs, &home.Rotograph{
		Title: heatestArticle.Title,
		Cover: heatestArticle.Cover,
		Color: "rgb(116,82,81)",
		Type:  "article",
		ToId:  heatestArticle.ID,
	})

	// 开启一个事务插入上面三条记录
	err = global.MysqlDb.Transaction(func(tx *gorm.DB) error {
		for _, rotograph := range rotographs {
			err := rotograph.Create(tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		global.Logger.Errorln("事务插入Rotograph出错err：" + err.Error())
	}
	return
}
