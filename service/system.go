package service

import (
	"time"
	"trojan/dao"
	"trojan/dao/redis"
	"trojan/module"
	"trojan/module/vo"
)

func SelectSystemByName(name *string) (*vo.SystemVo, error) {
	get := redis.Client.String.Get("trojan-panel:system")
	systemVo := new(vo.SystemVo)
	if err := get.ScanStruct(systemVo); err != nil {
		systemVo, err := dao.SelectSystemByName(name)
		if err != nil {
			return nil, err
		}
		redis.Client.String.Set("trojan-panel:system", systemVo, time.Hour.Milliseconds()*2)
		return systemVo, nil
	}
	return systemVo, nil
}

func UpdateSystemById(system *module.System) error {
	if err := dao.UpdateSystemById(system); err != nil {
		return err
	}
	redis.Client.Key.Del("trojan-panel:system")
	return nil
}
