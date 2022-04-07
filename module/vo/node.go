package vo

import "time"

type NodeVo struct {
	Id         uint      `json:"id"`
	Name       string    `json:"name"`
	Ip         string    `json:"ip"`
	Port       uint      `json:"port"`
	Type       uint      `json:"type"`
	CreateTime time.Time `json:"createTime"`
}

type NodePageVo struct {
	Nodes []NodeVo `json:"nodes"`
	BaseVoPage
}
