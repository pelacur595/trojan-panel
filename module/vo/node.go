package vo

import "time"

type NodeVo struct {
	Id         uint      `json:"id"`
	NodeSubId  uint      `json:"node_sub_id"`
	NodeTypeId uint      `json:"node_type_id"`
	Name       string    `json:"name"`
	Ip         string    `json:"ip"`
	Port       uint      `json:"port"`
	CreateTime time.Time `json:"createTime"`

	Ping int `json:"ping"`
}

type NodePageVo struct {
	Nodes []NodeVo `json:"nodes"`
	BaseVoPage
}
