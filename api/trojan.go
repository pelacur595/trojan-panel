package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"trojan/core"
	"trojan/module/vo"
)

// 停止trojan-gfw
func TrojanGFWStop(c *gin.Context) {
	if err := core.TrojanGFWStop(); err != nil {
		logrus.Errorln(err.Error())
		vo.Fail("停止trojan-gfw失败", c)
	} else {
		vo.Success("停止trojan-gfw成功", c)
	}
}

// 重启trojan-gfw
func TrojanGFWRestart(c *gin.Context) {
	if err := core.TrojanGFWRestart(); err != nil {
		logrus.Errorln(err.Error())
		vo.Fail("重启trojan-gfw失败!", c)
	} else {
		vo.Success("重启trojan-gfw成功!", c)
	}
}

// trojan-gfw状态
func TrojanGFWStatus(c *gin.Context) {
	vo.Success(core.TrojanGFWStatus(), c)
}

// 停止trojan-go
func TrojanGOStop(c *gin.Context) {
	if err := core.TrojanGOStop(); err != nil {
		logrus.Errorln(err.Error())
		vo.Fail("停止trojan-go失败", c)
	} else {
		vo.Success("停止trojan-go成功", c)
	}
}

// 重启trojan-go
func TrojanGORestart(c *gin.Context) {
	if err := core.TrojanGORestart(); err != nil {
		logrus.Errorln(err.Error())
		vo.Fail("重启trojan-go失败!", c)
	} else {
		vo.Success("重启trojan-go成功!", c)
	}
}

// trojan-go状态
func TrojanGOStatus(c *gin.Context) {
	vo.Success(core.TrojanGOStatus(), c)
}
