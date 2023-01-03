package service

import (
	"encoding/base64"
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

		systemAccountConfigBo := bo.SystemAccountConfigBo{}
		if err = json.Unmarshal([]byte(*system.AccountConfig), &systemAccountConfigBo); err != nil {
			logrus.Errorln(fmt.Sprintf("SelectSystemByName SystemAccountConfigBo 反序列化失败 err: %v", err))
			return systemVo, errors.New(constant.SysError)
		}
		systemEmailConfigBo := bo.SystemEmailConfigBo{}
		if err = json.Unmarshal([]byte(*system.EmailConfig), &systemEmailConfigBo); err != nil {
			logrus.Errorln(fmt.Sprintf("SelectSystemByName SystemEmailConfigBo 反序列化失败 err: %v", err))
			return systemVo, errors.New(constant.SysError)
		}
		systemTemplateConfigBo := bo.SystemTemplateConfigBo{}
		if err = json.Unmarshal([]byte(*system.TemplateConfig), &systemTemplateConfigBo); err != nil {
			logrus.Errorln(fmt.Sprintf("SelectSystemByName SystemTemplateConfigBo 反序列化失败 err: %v", err))
			return systemVo, errors.New(constant.SysError)
		}
		if systemTemplateConfigBo.ClashRule != "" {
			decodeString, err := base64.StdEncoding.DecodeString(systemTemplateConfigBo.ClashRule)
			if err != nil {
				logrus.Errorln(fmt.Sprintf("system config clashRule base64 stdEncoding decodeString err: %v", err))
				return systemVo, errors.New(constant.SysError)
			}
			systemTemplateConfigBo.ClashRule = string(decodeString)
		}

		systemVo = vo.SystemVo{
			Id:                          *system.Id,
			RegisterEnable:              systemAccountConfigBo.RegisterEnable,
			RegisterQuota:               systemAccountConfigBo.RegisterQuota,
			RegisterExpireDays:          systemAccountConfigBo.RegisterExpireDays,
			ResetDownloadAndUploadMonth: systemAccountConfigBo.ResetDownloadAndUploadMonth,
			TrafficRankEnable:           systemAccountConfigBo.TrafficRankEnable,
			ExpireWarnEnable:            systemEmailConfigBo.ExpireWarnEnable,
			ExpireWarnDay:               systemEmailConfigBo.ExpireWarnDay,
			EmailEnable:                 systemEmailConfigBo.EmailEnable,
			EmailHost:                   systemEmailConfigBo.EmailHost,
			EmailPort:                   systemEmailConfigBo.EmailPort,
			EmailUsername:               systemEmailConfigBo.EmailUsername,
			EmailPassword:               systemEmailConfigBo.EmailPassword,
			SystemName:                  systemTemplateConfigBo.SystemName,
			ClashRule:                   systemTemplateConfigBo.ClashRule,
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
	accountConfigBo := bo.SystemAccountConfigBo{}
	if systemDto.RegisterEnable != nil {
		accountConfigBo.RegisterEnable = *systemDto.RegisterEnable
	}
	if systemDto.RegisterQuota != nil {
		accountConfigBo.RegisterQuota = *systemDto.RegisterQuota
	}
	if systemDto.RegisterExpireDays != nil {
		accountConfigBo.RegisterExpireDays = *systemDto.RegisterExpireDays
	}
	if systemDto.ResetDownloadAndUploadMonth != nil {
		accountConfigBo.ResetDownloadAndUploadMonth = *systemDto.ResetDownloadAndUploadMonth
	}
	if systemDto.TrafficRankEnable != nil {
		accountConfigBo.TrafficRankEnable = *systemDto.TrafficRankEnable
	}
	accountConfigBoByte, err := json.Marshal(accountConfigBo)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("UpdateSystemById SystemAccountConfigBo 序列化异常err: %v", err))
	}
	accountConfigBoJsonStr := string(accountConfigBoByte)

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

	systemTemplateConfigBo := bo.SystemTemplateConfigBo{}
	if systemDto.SystemName != nil {
		systemTemplateConfigBo.SystemName = *systemDto.SystemName
	}
	if systemDto.ClashRule != nil {
		systemTemplateConfigBo.ClashRule = *systemDto.ClashRule
	}
	systemTemplateConfigBoByte, err := json.Marshal(systemTemplateConfigBo)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("UpdateSystemById SystemTemplateConfigBo 序列化异常err: %v", err))
	}
	systemTemplateConfigBoStr := string(systemTemplateConfigBoByte)

	system := module.System{
		Id:             systemDto.Id,
		AccountConfig:  &accountConfigBoJsonStr,
		EmailConfig:    &systemEmailConfigBoStr,
		TemplateConfig: &systemTemplateConfigBoStr,
	}

	if err := dao.UpdateSystemById(&system); err != nil {
		return err
	}
	redis.Client.Key.Del("trojan-panel:system")
	return nil
}
