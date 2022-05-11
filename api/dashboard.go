package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	redisgo "github.com/gomodule/redigo/redis"
	"trojan/dao/redis"
	"trojan/module/constant"
	"trojan/module/vo"
	"trojan/service"
)

func PanelGroup(c *gin.Context) {
	panelGroup, err := service.PanelGroup(c)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(panelGroup, c)
}

// 流量排行榜
func TrafficRank(c *gin.Context) {
	bytes, err := redis.Client.String.Get("trojan-panel:trafficRank").Bytes()
	if err != nil && err != redisgo.ErrNil {
		vo.Fail(err.Error(), c)
		return
	}
	if len(bytes) > 0 {
		var usersTrafficRankVo []vo.UsersTrafficRankVo
		if err := json.Unmarshal(bytes, &usersTrafficRankVo); err != nil {
			fmt.Println(err.Error())
			vo.Fail(constant.SysError, c)
			return
		}
		vo.Success(usersTrafficRankVo, c)
	} else {
		trafficRank, err := service.TrafficRank()
		if err != nil {
			vo.Fail(err.Error(), c)
			return
		}
		vo.Success(trafficRank, c)
	}
}
