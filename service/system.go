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
	"trojan/module/vo"
)

func SelectSystemByName(name *string) (*vo.SystemVo, error) {
	bytes, err := redis.Client.String.Get("trojan-panel:system").Bytes()
	if err != nil && err != redisgo.ErrNil {
		return nil, errors.New(constant.SysError)
	}
	if len(bytes) > 0 {
		var systemVo vo.SystemVo
		if err := json.Unmarshal(bytes, &systemVo); err != nil {
			logrus.Errorln(fmt.Sprintf("SystemVo JSON反转失败 err: %v", err))
			return nil, errors.New(constant.SysError)
		}
		return &systemVo, nil
	} else {
		systemVo, err := dao.SelectSystemByName(name)
		if err != nil {
			return nil, err
		}
		systemVoJson, err := json.Marshal(systemVo)
		if err != nil {
			logrus.Errorln(fmt.Sprintf("SystemVo JSON转换失败 err: %v", err))
			return nil, errors.New(constant.SysError)
		}
		redis.Client.String.Set("trojan-panel:system", systemVoJson, time.Minute.Microseconds()*30)
		return systemVo, nil
	}
}

func UpdateSystemById(system *module.System) error {
	if err := dao.UpdateSystemById(system); err != nil {
		return err
	}
	redis.Client.Key.Del("trojan-panel:system")
	return nil
}
