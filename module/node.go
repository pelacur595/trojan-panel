package module

import "time"

type Node struct {
	Id           *uint      `ddb:"id"`
	NodeServerId *uint      `ddb:"node_server_id"`
	NodeSubId    *uint      `ddb:"node_sub_id"`
	NodeTypeId   *uint      `ddb:"node_type_id"`
	Name         *string    `ddb:"name"`
	Ip           *string    `ddb:"ip"`
	Domain       *string    `ddb:"domain"`
	Port         *uint      `ddb:"port"`
	CreateTime   *time.Time `ddb:"create_time"`
	UpdateTime   *time.Time `ddb:"update_time"`
}
