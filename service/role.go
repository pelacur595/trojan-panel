package service

import (
	"trojan/dao"
	"trojan/module/dto"
	"trojan/module/vo"
)

func SelectRoleList(roleDto dto.RoleDto) (*[]vo.RoleListVo, error) {
	roleListVos, err := dao.SelectRoleList(roleDto)
	if err != nil {
		return nil, err
	}
	return roleListVos, nil
}
