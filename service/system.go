package service

import (
	"encoding/json"
	"errors"
	"fmt"
	redisgo "github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"os"
	"time"
	"trojan-panel/dao"
	"trojan-panel/dao/redis"
	"trojan-panel/model"
	"trojan-panel/model/bo"
	"trojan-panel/model/constant"
	"trojan-panel/model/dto"
	"trojan-panel/model/vo"
)

func SelectSystemByName(name *string) (vo.SystemVo, error) {
	var systemVo vo.SystemVo
	bytes, err := redis.Client.String.Get("trojan-panel:system").Bytes()
	if err != nil && err != redisgo.ErrNil {
		return systemVo, errors.New(constant.SysError)
	}
	if len(bytes) > 0 {
		if err = json.Unmarshal(bytes, &systemVo); err != nil {
			logrus.Errorln(fmt.Sprintf("SelectSystemByName SystemVo deserialization err: %v", err))
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
			logrus.Errorln(fmt.Sprintf("SelectSystemByName SystemAccountConfigBo deserialization err: %v", err))
			return systemVo, errors.New(constant.SysError)
		}
		systemEmailConfigBo := bo.SystemEmailConfigBo{}
		if err = json.Unmarshal([]byte(*system.EmailConfig), &systemEmailConfigBo); err != nil {
			logrus.Errorln(fmt.Sprintf("SelectSystemByName SystemEmailConfigBo deserialization err: %v", err))
			return systemVo, errors.New(constant.SysError)
		}
		systemTemplateConfigBo := bo.SystemTemplateConfigBo{}
		if err = json.Unmarshal([]byte(*system.TemplateConfig), &systemTemplateConfigBo); err != nil {
			logrus.Errorln(fmt.Sprintf("SelectSystemByName SystemTemplateConfigBo deserialization err: %v", err))
			return systemVo, errors.New(constant.SysError)
		}
		// 读取Clash规则默认模板文件
		clashRuleContent, err := os.ReadFile(constant.ClashRuleFilePath)
		if err != nil {
			logrus.Errorln(fmt.Sprintf("failed to read default template of Clash rule err: %v", err))
			return systemVo, errors.New(constant.SysError)
		}
		// 读取Xray默认模板文件
		xrayTemplateContent, err := os.ReadFile(constant.XrayTemplateFilePath)
		if err != nil {
			logrus.Errorln(fmt.Sprintf("failed to read Xray default template file err: %v", err))
			return systemVo, errors.New(constant.SysError)
		}

		systemVo = vo.SystemVo{
			Id:                          *system.Id,
			RegisterEnable:              systemAccountConfigBo.RegisterEnable,
			RegisterQuota:               systemAccountConfigBo.RegisterQuota,
			RegisterExpireDays:          systemAccountConfigBo.RegisterExpireDays,
			ResetDownloadAndUploadMonth: systemAccountConfigBo.ResetDownloadAndUploadMonth,
			TrafficRankEnable:           systemAccountConfigBo.TrafficRankEnable,
			CaptchaEnable:               systemAccountConfigBo.CaptchaEnable,
			ExpireWarnEnable:            systemEmailConfigBo.ExpireWarnEnable,
			ExpireWarnDay:               systemEmailConfigBo.ExpireWarnDay,
			EmailEnable:                 systemEmailConfigBo.EmailEnable,
			EmailHost:                   systemEmailConfigBo.EmailHost,
			EmailPort:                   systemEmailConfigBo.EmailPort,
			EmailUsername:               systemEmailConfigBo.EmailUsername,
			EmailPassword:               systemEmailConfigBo.EmailPassword,
			SystemName:                  systemTemplateConfigBo.SystemName,
			ClashRule:                   string(clashRuleContent),
			XrayTemplate:                string(xrayTemplateContent),
		}

		systemVoJson, err := json.Marshal(systemVo)
		if err != nil {
			logrus.Errorln(fmt.Sprintf("SelectSystemByName SystemVo serialization err: %v", err))
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
	if systemDto.CaptchaEnable != nil {
		accountConfigBo.CaptchaEnable = *systemDto.CaptchaEnable
	}
	accountConfigBoByte, err := json.Marshal(accountConfigBo)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("UpdateSystemById SystemAccountConfigBo serialization err: %v", err))
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
		logrus.Errorln(fmt.Sprintf("UpdateSystemById SystemEmailConfigBo serialization err: %v", err))
	}
	systemEmailConfigBoStr := string(systemEmailConfigBoByte)

	systemTemplateConfigBo := bo.SystemTemplateConfigBo{}
	if systemDto.SystemName != nil {
		systemTemplateConfigBo.SystemName = *systemDto.SystemName
	}
	if systemDto.ClashRule != nil {
		// 修改Clash规则默认模板文件
		if err := os.WriteFile(constant.ClashRuleFilePath, []byte(*systemDto.ClashRule), 0666); err != nil {
			logrus.Errorln(fmt.Sprintf("write Clash rule default template file err: %v", err))
		}
	}
	if systemDto.XrayTemplate != nil {
		// 修改Xray默认模板文件
		xrayConfigBo := bo.XrayConfigBo{}
		// 将json字符串映射到模板对象
		if err = json.Unmarshal([]byte(*systemDto.XrayTemplate), &xrayConfigBo); err != nil {
			logrus.Errorf("systemDto XrayTemplate deserialization err: %v", err)
			return err
		}
		xrayConfigBoStr, err := json.MarshalIndent(xrayConfigBo, "", "    ")
		if err != nil {
			logrus.Errorf("xrayConfigBo serialization err: %v", err)
			return err
		}
		if err := os.WriteFile(constant.XrayTemplateFilePath, xrayConfigBoStr, 0666); err != nil {
			logrus.Errorln(fmt.Sprintf("write Xray default template file err: %v", err))
		}
	}
	systemTemplateConfigBoByte, err := json.Marshal(systemTemplateConfigBo)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("UpdateSystemById SystemTemplateConfigBo serialization err: %v", err))
	}
	systemTemplateConfigBoStr := string(systemTemplateConfigBoByte)

	system := model.System{
		Id:             systemDto.Id,
		AccountConfig:  &accountConfigBoJsonStr,
		EmailConfig:    &systemEmailConfigBoStr,
		TemplateConfig: &systemTemplateConfigBoStr,
	}

	if err := dao.UpdateSystemById(&system); err != nil {
		return err
	}
	_ = redis.Client.Key.RetryDel("trojan-panel:system")
	return nil
}
