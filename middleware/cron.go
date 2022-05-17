package middleware

import (
	"github.com/robfig/cron/v3"
	"time"
	"trojan/service"
)

// 初始化定时任务
func InitCron() {
	location, _ := time.LoadLocation("Asia/Shanghai")
	c := cron.New(cron.WithLocation(location))
	// 扫描被禁用和到期的用户
	_, _ = c.AddFunc("@every 1m", service.ScanUsers)
	// 每天中午12点发送到期提醒邮件
	_, _ = c.AddFunc("0 0 12 * * *", service.ScanUserExpireWarn)
	// 每隔一小时刷新流量排行缓存
	_, _ = c.AddFunc("@every 1h", service.TrafficRankJob)
	c.Start()
}
