package service

import (
	"trojan/dao"
	"trojan/module"
	"trojan/module/dto"
)

func SelectRoleList(roleDto dto.RoleDto) ([]module.Role, error) {
	return dao.SelectRoleList(roleDto)
}
