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

	TrojanGoSni             *string `json:"trojanGoSni"`
	TrojanGoMuxEnable       *uint   `json:"trojanGoMuxEnable"`
	TrojanGoWebsocketEnable *uint   `json:"trojanGoWebsocketEnable"`
	TrojanGoWebsocketPath   *string `json:"trojanGoWebsocketPath"`
	TrojanGoWebsocketHost   *string `json:"trojanGoWebsocketHost"`
	TrojanGoSsEnable        *uint   `json:"trojanGoSsEnable"`
	TrojanGoSsMethod        *string `json:"trojanGoSsMethod"`
	TrojanGoSsPassword      *string `json:"trojanGoSsPassword"`
	HysteriaProtocol        *string `json:"hysteriaProtocol"`
	HysteriaUpMbps          *int    `json:"hysteriaUpMbps"`
	HysteriaDownMbps        *int    `json:"hysteriaDownMbps"`
	XrayProtocol            *string `json:"xrayProtocol"`
	XraySettings            *string `json:"xraySettings"`
	XrayStreamSettings      *string `json:"xrayStreamSettings"`
	XrayTag                 *string `json:"xrayTag"`
	XraySniffing            *string `json:"xraySniffing"`
	XrayAllocate            *string `json:"xrayAllocate"`
}

type NodePageVo struct {
	Nodes []NodeVo `json:"nodes"`
	BaseVoPage
}
