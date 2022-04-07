package module

import "time"

// 菜单
type MenuList struct {
	Id         *uint      `ddb:"id"`
	Name       *string    `ddb:"name"`
	Icon       *string    `ddb:"icon"`
	Route      *string    `ddb:"route"`
	Order      *uint      `ddb:"order"`
	ParentId   *uint      `ddb:"parent_id"`
	Path       *string    `ddb:"path"`
	Level      *uint      `ddb:"level"`
	CreateTime *time.Time `ddb:"create_time"`
	UpdateTime *time.Time `ddb:"update_time"`
}
