package service

import (
	"encoding/json"
	"errors"
	"fmt"
	redisgo "github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"time"
	"trojan-panel/dao"
	"trojan-panel/dao/redis"
	"trojan-panel/module"
	"trojan-panel/module/bo"
	"trojan-panel/module/constant"
	"trojan-panel/module/dto"
	"trojan-panel/module/vo"
)

func SelectSystemByName(name *string) (vo.SystemVo, error) {
	var systemVo vo.SystemVo
	bytes, err := redis.Client.String.Get("trojan-panel:system").Bytes()
	if err != nil && err != redisgo.ErrNil {
		return systemVo, errors.New(constant.SysError)
	}
	if len(bytes) > 0 {
		if err = json.Unmarshal(bytes, &systemVo); err != nil {
			logrus.Errorln(fmt.Sprintf("SelectSystemByName SystemVo 反序列化失败 err: %v", err))
			return systemVo, errors.New(constant.SysError)
		}
		return systemVo, nil
	} else {
		system, err := dao.SelectSystemByName(name)
		if err != nil {
			return systemVo, err
		}

		systemRegisterConfigBo := bo.SystemRegisterConfigBo{}
		if err = json.Unmarshal([]byte(*system.RegisterConfig), &systemRegisterConfigBo); err != nil {
			logrus.Errorln(fmt.Sprintf("SelectSystemByName SystemRegisterConfigBo 反序列化失败 err: %v", err))
			return systemVo, errors.New(constant.SysError)
		}
		systemEmailConfigBo := bo.SystemEmailConfigBo{}
		if err = json.Unmarshal([]byte(*system.EmailConfig), &systemEmailConfigBo); err != nil {
			logrus.Errorln(fmt.Sprintf("SelectSystemByName SystemEmailConfigBo 反序列化失败 err: %v", err))
			return systemVo, errors.New(constant.SysError)
		}

		systemVo = vo.SystemVo{
			Id:                 *system.Id,
			RegisterEnable:     systemRegisterConfigBo.RegisterEnable,
			RegisterQuota:      systemRegisterConfigBo.RegisterQuota,
			RegisterExpireDays: systemRegisterConfigBo.RegisterExpireDays,
			ResetQuotaMonth:    systemRegisterConfigBo.ResetQuotaMonth,
			TrafficRankEnable:  systemRegisterConfigBo.TrafficRankEnable,
			ExpireWarnEnable:   systemEmailConfigBo.ExpireWarnEnable,
			ExpireWarnDay:      systemEmailConfigBo.ExpireWarnDay,
			EmailEnable:        systemEmailConfigBo.EmailEnable,
			EmailHost:          systemEmailConfigBo.EmailHost,
			EmailPort:          systemEmailConfigBo.EmailPort,
			EmailUsername:      systemEmailConfigBo.EmailUsername,
			EmailPassword:      systemEmailConfigBo.EmailPassword,
		}

		systemVoJson, err := json.Marshal(systemVo)
		if err != nil {
			logrus.Errorln(fmt.Sprintf("SelectSystemByName SystemVo 序列化失败 err: %v", err))
			return systemVo, errors.New(constant.SysError)
		}
		redis.Client.String.Set("trojan-panel:system", systemVoJson, time.Minute.Milliseconds()*30/1000)

		return systemVo, nil
	}
}

func UpdateSystemById(systemDto dto.SystemUpdateDto) error {
	registerConfigBo := bo.SystemRegisterConfigBo{}
	if systemDto.RegisterEnable != nil {
		registerConfigBo.RegisterEnable = *systemDto.RegisterEnable
	}
	if systemDto.RegisterQuota != nil {
		registerConfigBo.RegisterQuota = *systemDto.RegisterQuota
	}
	if systemDto.RegisterExpireDays != nil {
		registerConfigBo.RegisterExpireDays = *systemDto.RegisterExpireDays
	}
	if systemDto.ResetQuotaMonth != nil {
		registerConfigBo.ResetQuotaMonth = *systemDto.ResetQuotaMonth
	}
	if systemDto.TrafficRankEnable != nil {
		registerConfigBo.TrafficRankEnable = *systemDto.TrafficRankEnable
	}
	registerConfigBoByte, err := json.Marshal(registerConfigBo)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("UpdateSystemById SystemRegisterConfigBo 序列化异常err: %v", err))
	}
	registerConfigBoJsonStr := string(registerConfigBoByte)

	systemEmailConfigBo := bo.SystemEmailConfigBo{}
	if systemDto.ExpireWarnEnable != nil {
		systemEmailConfigBo.ExpireWarnEnable = *systemDto.ExpireWarnEnable
	}
	if systemDto.ExpireWarnDay != nil {
		systemEmailConfigBo.ExpireWarnDay = *systemDto.ExpireWarnDay
	}
	if systemDto.EmailEnable != nil {
		systemEmailConfigBo.EmailEnable = *systemDto.EmailEnable
	}
	if systemDto.EmailHost != nil {
		systemEmailConfigBo.EmailHost = *systemDto.EmailHost
	}
	if systemDto.EmailPort != nil {
		systemEmailConfigBo.EmailPort = *systemDto.EmailPort
	}
	if systemDto.EmailUsername != nil {
		systemEmailConfigBo.EmailUsername = *systemDto.EmailUsername
	}
	if systemDto.EmailPassword != nil {
		systemEmailConfigBo.EmailPassword = *systemDto.EmailPassword
	}
	systemEmailConfigBoByte, err := json.Marshal(systemEmailConfigBo)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("UpdateSystemById SystemEmailConfigBo 序列化异常err: %v", err))
	}
	systemEmailConfigBoStr := string(systemEmailConfigBoByte)

	system := module.System{
		Id:             systemDto.Id,
		RegisterConfig: &registerConfigBoJsonStr,
		EmailConfig:    &systemEmailConfigBoStr,
	}

	if err := dao.UpdateSystemById(&system); err != nil {
		return err
	}
	redis.Client.Key.Del("trojan-panel:system")
	return nil
}
