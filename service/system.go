package service

import (
	"trojan/dao"
	"trojan/module"
	"trojan/module/vo"
)

func SelectSystemByName(name *string) (*vo.SystemVo, error) {
	systemVo, err := dao.SelectSystemByName(name)
	if err != nil {
		return nil, err
	}
	return systemVo, nil
}

func UpdateSystemById(system *module.System) error {
	if err := dao.UpdateSystemById(system); err != nil {
		return err
	}
	return nil
}
