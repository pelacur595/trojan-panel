package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"trojan/dao"
	"trojan/dao/redis"
	"trojan/module/vo"
)

func TrafficRankJob() {
	_, _ = TrafficRank()
}

func TrafficRank() ([]vo.UsersTrafficRankVo, error) {
	trafficRank, err := dao.TrafficRank()
	if err != nil {
		return nil, err
	}
	trafficRankJson, err := json.Marshal(trafficRank)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("json转换失败 err: %v", err))
	}
	redis.Client.String.Set("trojan-panel:trafficRank", trafficRankJson, time.Hour.Microseconds()*2/1000)
	return trafficRank, nil
}
