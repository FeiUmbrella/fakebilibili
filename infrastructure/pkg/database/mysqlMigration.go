package database

import (
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
	"fakebilibili/infrastructure/pkg/global"
	"log"
)

func init() {
	// 自动建表
	migration()
}

func migration() {
	err := global.MysqlDb.Set("gorm:table_options", "charset=utf8mb4").
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
