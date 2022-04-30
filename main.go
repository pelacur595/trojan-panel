package main

import (
	"github.com/gin-gonic/gin"
	"trojan/api"
	"trojan/core"
	"trojan/dao"
	"trojan/dao/redis"
	"trojan/middleware"
	"trojan/router"
	"trojan/util"
)

func main() {
	r := gin.Default()
	router.Router(r)
	_ = r.Run(":8081")
}

func init() {
	// 文件
	util.InitFile()
	// 全局配置
	core.InitConfig()
	// 日志
	middleware.InitLog()
	// 数据库
	dao.InitDB()
	// Redis
	redis.InitRedis()
	// 定时任务
	middleware.InitCron()
	// 限流
	middleware.InitRateLimiter()
	// 参数校验
	api.InitValidator()
}
