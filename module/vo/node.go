package vo

import "time"

type NodeVo struct {
	Id               uint      `json:"id"`
	Name             string    `json:"name"`
	Ip               string    `json:"ip"`
	Port             uint      `json:"port"`
	Sni              string    `json:"sni"`
	Type             uint      `json:"type"`
	WebsocketEnable  uint      `json:"websocketEnable"`
	WebsocketPath    string    `json:"websocketPath"`
	SsEnable         uint      `json:"ssEnable"`
	SsMethod         string    `json:"ssMethod"`
	SsPassword       string    `json:"ssPassword"`
	HysteriaProtocol string    `json:"hysteriaProtocol"`
	HysteriaUpMbps   int       `json:"hysteriaUpMbps"`
	HysteriaDownMbps int       `json:"hysteriaDownMbps"`
	CreateTime       time.Time `json:"createTime"`

	Ping   int `json:"ping"`
	OnLine int `json:"onLine"`
}

type NodePageVo struct {
	Nodes []NodeVo `json:"nodes"`
	BaseVoPage
}
