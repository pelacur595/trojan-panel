package vo

import "time"

type NodeVo struct {
	Id              uint      `json:"id"`
	Name            string    `json:"name"`
	Ip              string    `json:"ip"`
	Port            uint      `json:"port"`
	Type            uint      `json:"type"`
	WebsocketEnable uint      `json:"websocketEnable"`
	WebsocketPath   string    `json:"websocketPath"`
	SsEnable        uint      `json:"ssEnable"`
	SsMethod        string    `json:"ssMethod"`
	SsPassword      string    `json:"ssPassword"`
	CreateTime      time.Time `json:"createTime"`
}

type NodePageVo struct {
	Nodes []NodeVo `json:"nodes"`
	BaseVoPage
}
