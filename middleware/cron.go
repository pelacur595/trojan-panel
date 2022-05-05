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
	c.Start()
}
