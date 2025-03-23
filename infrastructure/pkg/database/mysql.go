package database

import (
	"context"
	"fakebilibili/infrastructure/config"
	"fakebilibili/infrastructure/model/contribution/article"
	"fakebilibili/infrastructure/model/contribution/video"
	"fakebilibili/infrastructure/model/contribution/video/barrage"
	"fakebilibili/infrastructure/model/contribution/video/comments"
	"fakebilibili/infrastructure/model/cron_events"
	"fakebilibili/infrastructure/model/home"
	"fakebilibili/infrastructure/model/sundry"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/model/user/attention"
	"fakebilibili/infrastructure/model/user/chat"
	"fakebilibili/infrastructure/model/user/checkIn"
	"fakebilibili/infrastructure/model/user/collect"
	"fakebilibili/infrastructure/model/user/favorites"
	"fakebilibili/infrastructure/model/user/liveInfo"
	"fakebilibili/infrastructure/model/user/notice"
	"fakebilibili/infrastructure/model/user/record"
	"fmt"
	"github.com/sethvargo/go-retry"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// 全局DB
var MysqlDb *gorm.DB

type MyWriter struct {
	log *logrus.Logger
}

type MysqlDB struct {
	DB *gorm.DB
}

// Printf 自定义结构体MyWriter实现gorm/logger.Writer 接口; 利用MyWriter.log将MySql的错误信息打印到日志文件中
func (m *MyWriter) Printf(format string, args ...interface{}) {
	m.log.Errorf(fmt.Sprintf(format, args...))
}

// 初始化全局DB，连接MySQL
func init() {
	// 数据库连接参数 Host post pw
	var mysqlConfig = config.Config.SqlConfig
	// sql 日志记录
	myLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel:                  logger.Silent, // log level
			IgnoreRecordNotFoundError: true,          // 忽略“记录未找到”错误
			Colorful:                  true,          // 禁止彩色打印
		},
	)
	// 创建一个以10s为单位的斐波那契(1, 2, 3, 5, 8...)*10s 的退避策略
	b := retry.NewFibonacci(10 * time.Second)
	// 创建一个空白上下文
	ctx := context.Background()
	// 当连MySQL失败时，会根据上面定义的斐波那契退避策略的等待时间进行重连
	if err := retry.Do(ctx, retry.WithMaxRetries(5, b), func(ctx context.Context) error {
		// 创建MySQL连接
		var err error
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", mysqlConfig.Username, mysqlConfig.Password, mysqlConfig.IP, mysqlConfig.Port, mysqlConfig.Database)
		MysqlDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: myLogger, // 传入自定义的数据库log格式
		})
		if err != nil {
			return err
		}
		if MysqlDb.Error != nil {
			return MysqlDb.Error
		}
		return nil
	}); err != nil {
		log.Fatalf("重连5次仍无法链接MySQL：%v \n", err)
	}

	// 自动建表
	migration()
}

func migration() {
	err := MysqlDb.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&user.User{},
			&attention.Attention{},
			&checkIn.CheckIn{},
			&collect.Collect{},
			&favorites.Favorites{},
			&liveInfo.LiveInfo{},
			&notice.Notice{},
			&notice.VideoInfo{},
			&notice.Article{},
			&record.Record{},
			&chat.ChatsListInfo{},
			&chat.Msg{},

			&sundry.Img{},
			&sundry.Upload{},
			&sundry.TranscodingTask{},

			&home.Rotograph{},

			&cron_events.CronEvent{},
			&cron_events.RuntimeLogEntry{},

			&article.ArticlesContribution{},
			&article.Classification{},
			&article.Comment{},
			&article.Likes{},

			&barrage.Barrage{},
			&barrage.VideoInfo{},
			&comments.Comment{},
			&comments.VideoInfo{},
			&video.Likes{},
			&video.VideosContribution{},
			&video.WatchRecord{},
		)
	if err != nil {
		log.Fatalf("AutoMigrate建表失败：%v", err)
	}
	return
}
