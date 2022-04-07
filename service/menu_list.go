package service

import (
	"trojan/dao"
	"trojan/module/vo"
	"trojan/util"
)

func SelectMenuListByRoleId(roleId *uint) ([]vo.TreeNode, error) {
	menuListVos, err := dao.SelectMenuListByRoleId(roleId)
	if err != nil {
		return nil, err
	}
	var treeNode []vo.TreeNode
	for _, v := range menuListVos {
		treeNode = append(treeNode, vo.TreeNode{
			Id:       v.Id,
			Name:     v.Name,
			Icon:     v.Icon,
			ParentId: v.ParentId,
			Route:    v.Route,
		})
	}
	rootNode := util.BuildTree(0, treeNode)
	return rootNode, nil
}
