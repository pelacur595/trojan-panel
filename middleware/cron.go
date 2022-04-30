package middleware

import (
	"github.com/robfig/cron/v3"
	"time"
	"trojan/service"
)

// 初始化定时任务
func init() {
	location, _ := time.LoadLocation("Asia/Shanghai")
	c := cron.New(cron.WithLocation(location))
	_, _ = c.AddFunc("@every 5m", service.ScanUsers)
	c.Start()
}
