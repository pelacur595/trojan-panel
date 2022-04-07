package module

import "time"

// 角色菜单关系
type RoleMenuList struct {
	Id         *uint      `ddb:"id"`
	RoleId     *uint      `ddb:"role_id"`
	MenuListId *uint      `ddb:"menu_list_id"`
	CreateTime *time.Time `ddb:"create_time"`
	UpdateTime *time.Time `ddb:"update_time"`
}
