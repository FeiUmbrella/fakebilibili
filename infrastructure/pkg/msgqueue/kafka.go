package msgqueue

import (
	"context"
	"fakebilibili/infrastructure/config"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

var kafkaConfig = config.Config.KafkaConfig

// KafkaProducerPool todo：生产者链接池？
type KafkaProducerPool struct {
	ConnPool chan *kafka.Conn
}

func ReturnNormalInstance() *kafka.Conn {
	conn, err := GetNormalProducerConn()
	if err != nil {
		log.Fatalf("获取即时kafka连接出错：%v", err)
	}
	return conn
}

func ReturnDelayInstance() *kafka.Conn {
	conn, err := GetDelayProducerConn()
	if err != nil {
		log.Fatalf("获取延迟kafka连接出错：%v", err)
	}
	return conn
}

func GetNormalProducerConn() (*kafka.Conn, error) {
	// topic: normalTopic    partition:0
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaConfig.Server, kafkaConfig.NormalTopic, 0)
	if err != nil {
		return nil, fmt.Errorf("创建即时kafka连接出错%v", err)
	}
	return conn, nil
}

func GetDelayProducerConn() (*kafka.Conn, error) {
	// topic: delayTopic    partition:0
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaConfig.Server, kafkaConfig.DelayTopic, 0)
	if err != nil {
		return nil, fmt.Errorf("创建延迟kafka连接出错%v", err)
	}
	return conn, nil
}
