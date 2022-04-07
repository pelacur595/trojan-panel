package vo

type MenuListVo struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Route    string `json:"route"`
	ParentId uint   `json:"parentId"`
}

type TreeNode struct {
	Id       uint       `json:"id"`
	Name     string     `json:"name"`
	Icon     string     `json:"icon"`
	ParentId uint       `json:"parentId"`
	Route    string     `json:"route"`
	Children []TreeNode `json:"children"`
}
