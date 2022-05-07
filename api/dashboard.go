package api

import (
	"github.com/gin-gonic/gin"
	"trojan/core"
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
