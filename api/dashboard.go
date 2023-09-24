package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	redisgo "github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"trojan-panel/dao/redis"
	"trojan-panel/model/constant"
	"trojan-panel/model/vo"
	"trojan-panel/service"
)

func PanelGroup(c *gin.Context) {
	panelGroup, err := service.PanelGroup(c)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(panelGroup, c)
}

// TrafficRank 流量排行榜
func TrafficRank(c *gin.Context) {
	bytes, err := redis.Client.String.Get("trojan-panel:trafficRank").Bytes()
	if err != nil && err != redisgo.ErrNil {
		vo.Fail(constant.SysError, c)
		return
	}
	if len(bytes) > 0 {
		var accountTrafficRankVo []vo.AccountTrafficRankVo
		if err := json.Unmarshal(bytes, &accountTrafficRankVo); err != nil {
			logrus.Errorln(fmt.Sprintf("AccountTrafficRankVo deserialization err: %v", err))
			vo.Fail(constant.SysError, c)
			return
		}
		vo.Success(accountTrafficRankVo, c)
	} else {
		trafficRank, err := service.TrafficRank()
		if err != nil {
			vo.Fail(err.Error(), c)
			return
		}
		vo.Success(trafficRank, c)
	}
}
