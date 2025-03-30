package checkin

import (
	"errors"
	"fakebilibili/adapter/http/receive"
	"fakebilibili/infrastructure/model/user/checkIn"
	"fakebilibili/infrastructure/pkg/global"
)

// GetUserIntegral 获取用户积分
func GetUserIntegral(data *receive.GetUserIntegralRequest) (results interface{}, err error) {
	c := &checkIn.CheckIn{}
	if err := global.MysqlDb.Model(&checkIn.CheckIn{}).Where("uid = ?", data.UID).Find(&c).Error; err != nil {
		return nil, errors.New("查询积分出错")
	}
	return c.Integral, nil
}
