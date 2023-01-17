package service

import (
	"encoding/json"
	"errors"
	"fmt"
	redisgo "github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
	"trojan-panel/core"
	"trojan-panel/dao"
	"trojan-panel/dao/redis"
	"trojan-panel/module/bo"
	"trojan-panel/module/constant"
)

func SelectXrayTemplate() (bo.XrayTemplate, error) {
	var xrayTemplate bo.XrayTemplate
	bytes, err := redis.Client.String.Get("trojan-panel:config:template-xray").Bytes()
	if err != nil && err != redisgo.ErrNil {
		return xrayTemplate, errors.New(constant.SysError)
	}
	if len(bytes) > 0 {
		if err = json.Unmarshal(bytes, &xrayTemplate); err != nil {
			logrus.Errorln(fmt.Sprintf("SelectXrayTemplate XrayTemplate 反序列化失败 err: %v", err))
			return xrayTemplate, errors.New(constant.SysError)
		}
		return xrayTemplate, nil
	} else {
		xrayTemplateContent, err := os.ReadFile(constant.XrayTemplateFilePath)
		if err != nil {
			logrus.Errorln(fmt.Sprintf("读取Xray模板失败 err: %v", err))
			return xrayTemplate, errors.New(constant.SysError)
		}
		xrayTemplateJson, err := json.Marshal(xrayTemplateContent)
		if err != nil {
			logrus.Errorln(fmt.Sprintf("SelectXrayTemplate XrayTemplate 序列化失败 err: %v", err))
			return xrayTemplate, errors.New(constant.SysError)
		}
		redis.Client.String.Set("trojan-panel:config:template-xray", xrayTemplateJson, time.Minute.Milliseconds()*30/1000)
		return xrayTemplate, nil
	}
}

func UpdateXrayTemplate(token string, xrayTemplate string) error {
	var mutex sync.Mutex
	defer mutex.Unlock()
	if mutex.TryLock() {
		ips, err := dao.SelectNodesIpDistinct()
		if err != nil {
			return err
		}
		for _, ip := range ips {
			xrayTemplateDto := core.XrayTemplateDto{XrayTemplate: xrayTemplate}
			if err := core.UpdateXrayTemplate(token, ip, &xrayTemplateDto); err != nil {
				logrus.Errorln(fmt.Sprintf("更新Xray默认模板异常 ip: %s err: %v", ip, err))
			}
		}
		if err := os.WriteFile(constant.XrayTemplateFilePath, []byte(xrayTemplate), 0666); err != nil {
			logrus.Errorln(fmt.Sprintf("写入Xray默认模板异常err: %v", err))
			return errors.New(constant.SysError)
		}
		redis.Client.Key.Del("trojan-panel:config:template-xray")
	}
	return nil
}

func SyncXrayTemplate(token string) {
	ips, err := dao.SelectNodesIpDistinct()
	if err != nil {
		logrus.Errorln(fmt.Sprintf("查询服务器IP异常 err: %v", err))
		return
	}
	xrayTemplate, err := SelectXrayTemplate()
	if err != nil {
		logrus.Errorln(fmt.Sprintf("查询Xray默认模板异常 err: %v", err))
		return
	}
	xrayTemplateJson, err := json.Marshal(xrayTemplate)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Xray默认模板序列化异常 err: %v", err))
		return
	}
	for _, ip := range ips {
		xrayTemplateDto := core.XrayTemplateDto{XrayTemplate: string(xrayTemplateJson)}
		if err := core.UpdateXrayTemplate(token, ip, &xrayTemplateDto); err != nil {
			logrus.Errorln(fmt.Sprintf("更新Xray默认模板异常 ip: %s err: %v", ip, err))
		}
	}
}
