package service

import (
	"errors"
	"trojan/dao"
	"trojan/dao/redis"
	"trojan/module"
	"trojan/module/constant"
	"trojan/module/vo"
)

func SelectSystemByName(name *string) (*vo.SystemVo, error) {
	get := redis.Client.Hash.HGetAll("trojan-panel:system")
	values, err := get.Values()
	if err != nil {
		return nil, errors.New(constant.SysError)
	}
	if len(values) == 0 {
		systemVo, err := dao.SelectSystemByName(name)
		if err != nil {
			return nil, err
		}
		redis.Client.Hash.HMSetFromStruct("trojan-panel:system", systemVo)
		return systemVo, nil
	} else {
		systemVo := vo.SystemVo{}
		if err := get.ScanStruct(&systemVo); err != nil {
			return nil, errors.New(constant.SysError)
		}
		return &systemVo, nil
	}
}

func UpdateSystemById(system *module.System) error {
	if err := dao.UpdateSystemById(system); err != nil {
		return err
	}
	redis.Client.Key.Del("trojan-panel:system")
	return nil
}
