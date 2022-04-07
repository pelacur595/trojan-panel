package util

import "trojan/module/vo"

// 构建树形结构
func BuildTree(parentId uint, data []vo.TreeNode) []vo.TreeNode {
	var treeNode []vo.TreeNode
	for _, v := range data {
		if v.ParentId == parentId {
			var children []vo.TreeNode
			children = append(children, BuildTree(v.Id, data)...)
			treeNode = append(treeNode, vo.TreeNode{
				Id:       v.Id,
				Name:     v.Name,
				Icon:     v.Icon,
				ParentId: v.ParentId,
				Route:    v.Route,
				Children: children,
			})
		}
	}
	return treeNode
}
