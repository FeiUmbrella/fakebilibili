package global

import (
	"fakebilibili/infrastructure/config"
	"fakebilibili/infrastructure/pkg/database/mysql"
	redis2 "fakebilibili/infrastructure/pkg/database/redis"
	log "fakebilibili/infrastructure/pkg/logrus"
	"fakebilibili/infrastructure/pkg/msgqueue"
	"github.com/go-redis/redis"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	Config         *config.Info
	Logger         *logrus.Logger
	MysqlDb        *gorm.DB
	RedisDb        *redis.Client
	NormalProducer *kafka.Conn
	DelayProducer  *kafka.Conn
)

// init 该包被导入时自动执行，将各个全局变量汇聚在一起
func init() {
	Config = config.Config   // 全局配置参数
	Logger = log.Logger      // 全局日志
	MysqlDb = mysql.MysqlDb  // 全局数据库
	RedisDb = redis2.RedisDb // 全局Redis
	NormalProducer = msgqueue.ReturnNormalInstance()
	DelayProducer = msgqueue.ReturnDelayInstance()
}
