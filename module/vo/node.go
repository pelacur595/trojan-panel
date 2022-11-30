package vo

import "time"

// NodeVo 查询分页Node对象
type NodeVo struct {
	Id         uint      `json:"id"`
	NodeSubId  uint      `json:"nodeSubId"`
	NodeTypeId uint      `json:"nodeTypeId"`
	Name       string    `json:"name"`
	Ip         string    `json:"ip"`
	Port       uint      `json:"port"`
	CreateTime time.Time `json:"createTime"`
}

type NodePageVo struct {
	Nodes []NodeVo `json:"nodes"`
	BaseVoPage
}

// NodeOneVo 查询单个Node对象
type NodeOneVo struct {
	Id         uint      `json:"id"`
	NodeSubId  uint      `json:"nodeSubId"`
	NodeTypeId uint      `json:"nodeTypeId"`
	Name       string    `json:"name"`
	Ip         string    `json:"ip"`
	Port       uint      `json:"port"`
	CreateTime time.Time `json:"createTime"`

	XrayProtocol             string                   `json:"xrayProtocol"`
	XraySettings             string                   `json:"xraySettings"`
	XrayStreamSettingsEntity XrayStreamSettingsEntity `json:"xrayStreamSettingsEntity"`
	XrayTag                  string                   `json:"xrayTag"`
	XraySniffing             string                   `json:"xraySniffing"`
	XrayAllocate             string                   `json:"xrayAllocate"`
	TrojanGoSni              string                   `json:"trojanGoSni"`
	TrojanGoMuxEnable        uint                     `json:"trojanGoMuxEnable"`
	TrojanGoWebsocketEnable  uint                     `json:"trojanGoWebsocketEnable"`
	TrojanGoWebsocketPath    string                   `json:"trojanGoWebsocketPath"`
	TrojanGoWebsocketHost    string                   `json:"trojanGoWebsocketHost"`
	TrojanGoSsEnable         uint                     `json:"trojanGoSsEnable"`
	TrojanGoSsMethod         string                   `json:"trojanGoSsMethod"`
	TrojanGoSsPassword       string                   `json:"trojanGoSsPassword"`
	HysteriaProtocol         string                   `json:"hysteriaProtocol"`
	HysteriaUpMbps           int                      `json:"hysteriaUpMbps"`
	HysteriaDownMbps         int                      `json:"hysteriaDownMbps"`
}

type XrayStreamSettingsEntity struct {
	Network    string                             `json:"network"`
	Security   string                             `json:"security"`
	WsSettings XrayStreamSettingsWsSettingsEntity `json:"wsSettings"`
}

type XrayStreamSettingsWsSettingsEntity struct {
	Path string `json:"path"`
}
