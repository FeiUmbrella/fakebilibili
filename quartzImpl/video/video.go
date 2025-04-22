package video

import (
	"fakebilibili/infrastructure/pkg/global"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

// PublishVideoOnSchedule 将定时视频写入到延迟队列
// todo:将预定发布时间写入延时队列里，再由演示队列的消费者去处理;为什么1次向延时队列提交了两条延时消息？？？？？而且其中的一条被立马转到了及时队列中去
func PublishVideoOnSchedule(scheduleTime time.Time, id uint) error {
	msg := kafka.Message{Value: []byte(fmt.Sprintf("publishVideo_%d", id)), Time: scheduleTime}
	_, err := global.DelayProducer.WriteMessages(msg)
	global.Logger.Infof("定时任务写入延时队列成功:%v", msg)
	if err != nil {
		global.Logger.Errorf("定时发布视频任务写入消息队列失败：%v", err)
		return err
	}
	return nil
}
