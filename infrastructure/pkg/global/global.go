package global

import (
	"fakebilibili/infrastructure/config"
	"fakebilibili/infrastructure/pkg/database"
	log "fakebilibili/infrastructure/pkg/logrus"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	Config  *config.Info
	Logger  *logrus.Logger
	MysqlDb *gorm.DB
	RedisDb *redis.Client
)

// init 该包被导入时自动执行，将各个全局变量汇聚在一起
func init() {
	Config = config.Config     // 全局配置参数
	Logger = log.Logger        // 全局日志
	MysqlDb = database.MysqlDb // 全局数据库
	RedisDb = database.RedisDb // 全局Redis
}
