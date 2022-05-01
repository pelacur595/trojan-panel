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
	util.InitFile()
	core.InitConfig()
	middleware.InitLog()
	dao.InitMySQL()
	redis.InitRedis()
	middleware.InitCron()
	middleware.InitRateLimiter()
	api.InitValidator()
}
