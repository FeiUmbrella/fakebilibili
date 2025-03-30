package checkin

import (
	"errors"
	"fakebilibili/adapter/http/receive"
	"fakebilibili/infrastructure/model/user/checkIn"
	"fakebilibili/infrastructure/pkg/global"
	"fakebilibili/infrastructure/pkg/utils/date"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// CheckIn 签到
func CheckIn(data *receive.CheckInRequestStruct) (interface{}, error) {
	check := &checkIn.CheckIn{Uid: data.UID}
	err := check.Query()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("查询签到历史出错")
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建第一次签到
		check.Uid = data.UID
		check.LatestDay = date.GetDay(time.Now())
		check.ConsecutiveDays = 1
		check.Integral = 1
		if !check.Create() {
			return nil, errors.New("创建签到记录失败")
		}
	} else { // 不是第一次签到
		if check.LatestDay == date.GetDay(time.Now()) {
			return "今天已签到，请勿重复签到", nil
		}
		if check.LatestDay == date.GetYesterday() { // 连续签到
			check.ConsecutiveDays += 1
		} else {
			// 非连续签到
			check.ConsecutiveDays = 1
		}
		check.Integral += 1 // 积分

		// 更新该用户签到记录
		if err := check.Update(map[string]interface{}{
			"consecutive_days": check.ConsecutiveDays,
			"integral":         check.Integral,
			"latest_day":       date.GetDay(time.Now()),
		}); err != nil {
			return nil, err
		}
	}
	global.Logger.Infof("用户%d成功签到，已经连续签到%d天", check.Uid, check.ConsecutiveDays)
	return fmt.Sprintf("签到成功，您已连续签到%d天", check.ConsecutiveDays), nil
}
