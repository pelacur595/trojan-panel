package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	redisgo "github.com/gomodule/redigo/redis"
	"trojan/core"
	"trojan/dao/redis"
	"trojan/module/constant"
	"trojan/module/vo"
	"trojan/service"
	"trojan/util"
)

func PanelGroup(c *gin.Context) {
	userInfo, err := service.GetUserInfo(c)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	userVo, err := service.SelectUserById(&userInfo.Id)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	nodeCount, err := service.CountNode()
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	panelGroupVo := vo.PanelGroupVo{
		TotalFlow:    userVo.Quota,
		ResidualFlow: userVo.Quota - int(userVo.Upload) - int(userVo.Download),
		NodeCount:    nodeCount,
		ExpireTime:   userVo.ExpireTime,
	}
	if util.IsAdmin(userInfo.Roles) {
		userCount, err := service.CountUserByUsername(nil)
		if err != nil {
			vo.Fail(err.Error(), c)
			return
		}
		panelGroupVo.UserCount = userCount

		// 在线用户
		api := core.TrojanGoApi()
		ips, err := service.SelectNodeIps()
		if err != nil {
			return
		}
		var online = 0
		for _, ip := range ips {
			num, err := api.OnLine(ip)
			if err != nil {
				return
			}
			online += num
		}
		panelGroupVo.OnLine = online
	}
	vo.Success(panelGroupVo, c)
}

// 流量排行榜
func TrafficRank(c *gin.Context) {
	bytes, err := redis.Client.String.Get("trojan-panel:trafficRank").Bytes()
	if err != nil && err != redisgo.ErrNil {
		vo.Fail(err.Error(), c)
		return
	}
	if len(bytes) > 0 {
		usersTrafficRankVo := vo.UsersTrafficRankVo{}
		if err := json.Unmarshal(bytes, &usersTrafficRankVo); err != nil {
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
