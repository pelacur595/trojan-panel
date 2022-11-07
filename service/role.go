package service

import (
	"trojan-panel/dao"
	"trojan-panel/module"
	"trojan-panel/module/dto"
)

func SelectRoleList(roleDto dto.RoleDto) ([]module.Role, error) {
	return dao.SelectRoleList(roleDto)
}
