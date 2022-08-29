package service

import (
	"encoding/json"
	"errors"
	"fmt"
	redisgo "github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"time"
	"trojan/dao"
	"trojan/dao/redis"
	"trojan/module"
	"trojan/module/constant"
)

func SelectSystemByName(name *string) (*module.System, error) {
	bytes, err := redis.Client.String.Get("trojan-panel:system").Bytes()
	if err != nil && err != redisgo.ErrNil {
		return nil, errors.New(constant.SysError)
	}
	if len(bytes) > 0 {
		var system module.System
		if err = json.Unmarshal(bytes, &system); err != nil {
			logrus.Errorln(fmt.Sprintf("SystemVo JSON反转失败 err: %v", err))
			return nil, errors.New(constant.SysError)
		}
		return &system, nil
	} else {
		system, err := dao.SelectSystemByName(name)
		if err != nil {
			return nil, err
		}
		systemJson, err := json.Marshal(system)
		if err != nil {
			logrus.Errorln(fmt.Sprintf("SystemVo JSON转换失败 err: %v", err))
			return nil, errors.New(constant.SysError)
		}
		redis.Client.String.Set("trojan-panel:system", systemJson, time.Minute.Milliseconds()*30/1000)
		return system, nil
	}
}

func UpdateSystemById(system *module.System) error {
	if err := dao.UpdateSystemById(system); err != nil {
		return err
	}
	redis.Client.Key.Del("trojan-panel:system")
	return nil
}
