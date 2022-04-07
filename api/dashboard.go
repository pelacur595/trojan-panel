package api

import (
	"github.com/gin-gonic/gin"
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
	if util.IsAdmin(userInfo.RoleNames) {
		userCount, err := service.CountUserByUsername(nil)
		if err != nil {
			vo.Fail(err.Error(), c)
			return
		}
		panelGroupVo.UserCount = userCount
	}
	vo.Success(panelGroupVo, c)
}
