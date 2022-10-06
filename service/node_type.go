package service

import (
	"trojan-panel/dao"
	"trojan-panel/module/vo"
)

func SelectNodeTypeList() ([]vo.NodeTypeVo, error) {
	return dao.SelectNodeTypeList()
}
