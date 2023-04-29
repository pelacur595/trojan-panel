package main

import (
	"github.com/gin-gonic/gin"
	"trojan-panel/api"
	"trojan-panel/core"
	"trojan-panel/dao"
	"trojan-panel/dao/redis"
	"trojan-panel/middleware"
	"trojan-panel/router"
	"trojan-panel/util"
)

func main() {
	r := gin.Default()
	router.Router(r)
	_ = r.Run(":8081")
	defer releaseResource()
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

func releaseResource() {
	dao.CloseDb()
	redis.CloseRedis()
}
