package util

import "time"

// 过期天数转过期时间（13位时间戳）
func DayToMilli(day uint) uint {
	return NowMilli() + day*86400*1000
}

// 获取当前时间的13位时间戳
func NowMilli() uint {
	return uint(time.Now().UnixMilli())
}
