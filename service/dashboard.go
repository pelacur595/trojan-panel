package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
	"trojan-panel/dao"
	"trojan-panel/dao/redis"
	"trojan-panel/module/constant"
	"trojan-panel/module/vo"
	"trojan-panel/util"
)

// CronTrafficRank 流量排行榜 一小时更新一次
func CronTrafficRank() {
	_, _ = TrafficRank()
}

func TrafficRank() ([]vo.AccountTrafficRankVo, error) {
	roleIds := []uint{constant.USER}
	trafficRank, err := dao.TrafficRank(&roleIds)
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
		logrus.Errorln(fmt.Sprintf("AccountTrafficRankVo JSON转换失败 err: %v", err))
		return nil, errors.New(constant.SysError)
	}
	redis.Client.String.Set("trojan-panel:trafficRank", trafficRankJson, time.Hour.Milliseconds()*2/1000)
	return trafficRank, nil
}

func PanelGroup(c *gin.Context) (*vo.PanelGroupVo, error) {
	accountInfo, err := GetAccountInfo(c)
	if err != nil {
		return nil, err
	}
	account, err := SelectAccountById(&accountInfo.Id)
	if err != nil {
		return nil, err
	}
	nodeCount, err := CountNode()
	if err != nil {
		return nil, err
	}
	panelGroupVo := vo.PanelGroupVo{
		Quota:        *account.Quota,
		ResidualFlow: *account.Quota - *account.Upload - *account.Download,
		NodeCount:    nodeCount,
		ExpireTime:   *account.ExpireTime,
	}
	if util.IsAdmin(accountInfo.Roles) {
		var err error
		accountCount, err := CountAccountByUsername(nil)
		cpuUsed, err := GetCpuPercent()
		memUsed, err := GetMemPercent()
		diskUsed, err := GetDiskPercent()
		if err != nil {
			return nil, err
		}
		panelGroupVo.AccountCount = accountCount
		panelGroupVo.CpuUsed = cpuUsed
		panelGroupVo.MemUsed = memUsed
		panelGroupVo.DiskUsed = diskUsed
	}
	return &panelGroupVo, nil
}

// GetCpuPercent 获取CPU使用率
func GetCpuPercent() (float64, error) {
	var err error
	percent, err := cpu.Percent(time.Second, false)
	value, err := strconv.ParseFloat(fmt.Sprintf("%.1f", percent[0]), 64)
	return value, err
}

// GetMemPercent 获取内存使用率
func GetMemPercent() (float64, error) {
	var err error
	memInfo, err := mem.VirtualMemory()
	value, err := strconv.ParseFloat(fmt.Sprintf("%.1f", memInfo.UsedPercent), 64)
	return value, err
}

// GetDiskPercent 获取硬盘使用率
func GetDiskPercent() (float64, error) {
	var err error
	parts, err := disk.Partitions(true)
	diskInfo, err := disk.Usage(parts[0].Mountpoint)
	value, err := strconv.ParseFloat(fmt.Sprintf("%.1f", diskInfo.UsedPercent), 64)
	return value, err
}
