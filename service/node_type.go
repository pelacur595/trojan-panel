package service

import (
	"trojan/dao"
	"trojan/module/vo"
)

func SelectNodeTypeList() ([]vo.NodeTypeVo, error) {
	return dao.SelectNodeTypeList()
}
