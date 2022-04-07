package service

import (
	"trojan/dao"
	"trojan/module/vo"
)

func SelectNodeTypeList() ([]vo.NodeTypeVo, error) {
	nodeTypeVos, err := dao.SelectNodeTypeList()
	if err != nil {
		return nil, err
	}
	return nodeTypeVos, nil
}
