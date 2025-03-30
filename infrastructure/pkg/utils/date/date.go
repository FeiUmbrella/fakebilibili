package date

import (
	"strconv"
	"time"
)

var TimeDay = "20060102"

// GetDay 获得传入时间的day
func GetDay(now time.Time) int {
	ret, _ := strconv.Atoi(now.Format(TimeDay))
	return ret
}

// GetYesterday 获得当前时间的昨天
func GetYesterday() int {
	ret, _ := strconv.Atoi(time.Now().Add(-1 * 24 * time.Hour).Format(TimeDay))
	return ret
}
