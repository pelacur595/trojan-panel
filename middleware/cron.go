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
	_, _ = c.AddFunc("@every 1m", service.ScanUsers)
	_, _ = c.AddFunc("0 0 12 * * *", service.ScanUserExpireWarn)
	_, _ = c.AddFunc("@every 1h", service.TrafficRankJob)
	c.Start()
}
