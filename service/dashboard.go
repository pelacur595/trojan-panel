package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
	"trojan/dao"
	"trojan/dao/redis"
	"trojan/module/constant"
	"trojan/module/vo"
	"trojan/util"
)

// 流量排行榜 一小时更新一次
func TrafficRankJob() {
	_, _ = TrafficRank()
}

func TrafficRank() ([]vo.UsersTrafficRankVo, error) {
	trafficRank, err := dao.TrafficRank()
	for index, item := range trafficRank {
		usernameLen := len(item.Username)
		prefix := item.Username[0:2]
		suffix := item.Username[usernameLen-2:]
		trafficRank[index].Username = fmt.Sprintf("%s****%s", prefix, suffix)
	}
	if err != nil {
		return nil, err
	}
	trafficRankJson, err := json.Marshal(trafficRank)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("UsersTrafficRankVo JSON转换失败 err: %v", err))
		return nil, errors.New(constant.SysError)
	}
	redis.Client.String.Set("trojan-panel:trafficRank", trafficRankJson, time.Hour.Milliseconds()*2/1000)
	return trafficRank, nil
}

func PanelGroup(c *gin.Context) (*vo.PanelGroupVo, error) {
	userInfo, err := GetUserInfo(c)
	if err != nil {
		return nil, err
	}
	userVo, err := SelectUserById(&userInfo.Id)
	if err != nil {
		return nil, err
	}
	nodeCount, err := CountNode()
	if err != nil {
		return nil, err
	}
	panelGroupVo := vo.PanelGroupVo{
		Quota:        userVo.Quota,
		ResidualFlow: userVo.Quota - int(userVo.Upload) - int(userVo.Download),
		NodeCount:    nodeCount,
		ExpireTime:   userVo.ExpireTime,
	}
	if util.IsAdmin(userInfo.Roles) {
		userCount, err := CountUserByUsername(nil)
		if err != nil {
			return nil, err
		}
		panelGroupVo.UserCount = userCount

		//// 在线用户
		//api := core.TrojanGoApi()
		//ips, err := SelectNodeIps()
		//if err != nil {
		//	return nil, err
		//}
		//var online = 0
		//for _, ip := range ips {
		//	num, err := api.OnLine(ip)
		//	if err != nil {
		//		continue
		//	}
		//	online += num
		//}
		//panelGroupVo.OnLine = online
	}
	return &panelGroupVo, nil
}
