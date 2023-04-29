package dto

type NodePageDto struct {
	NodeDto
	BaseDto
}

type NodeDto struct {
	Name         *string `json:"name" form:"name" validate:"omitempty,min=0,max=20"`
	NodeServerId *uint   `json:"nodeServerId" form:"nodeServerId" validate:"omitempty,gte=0"`
}

type NodeCreateDto struct {
	NodeServerId *uint   `json:"nodeServerId" form:"nodeServerId" validate:"required,gt=0"`
	NodeTypeId   *uint   `json:"nodeTypeId" form:"nodeTypeId" validate:"required,gt=0"`
	Name         *string `json:"name" form:"name" validate:"required,min=2,max=20"`
	Domain       *string `json:"domain" form:"domain" validate:"required,ip|fqdn,min=4,max=64"`
	Port         *uint   `json:"port" form:"port" validate:"required,validatePort"`

	XrayFlow           *string `json:"xrayFlow" form:"xrayFlow" validate:"omitempty,min=0,max=32"`
	XraySSMethod       *string `json:"xraySSMethod" form:"xraySSMethod" validate:"omitempty,min=0,max=32"`
	XrayProtocol       *string `json:"xrayProtocol" form:"xrayProtocol" validate:"omitempty,min=0,max=32"`
	XraySettings       *string `json:"xraySettings" form:"xraySettings" validate:"omitempty,min=0,max=256"`
	XrayStreamSettings *string `json:"xrayStreamSettings" form:"xrayStreamSettings" validate:"omitempty,min=0,max=256"`
	XrayTag            *string `json:"xrayTag" form:"xrayTag" validate:"omitempty,min=0,max=64"`
	XraySniffing       *string `json:"xraySniffing" form:"xraySniffing" validate:"omitempty,min=0,max=256"`
	XrayAllocate       *string `json:"xrayAllocate" form:"xrayAllocate" validate:"omitempty,min=0,max=256"`

	TrojanGoSni             *string `json:"trojanGoSni" form:"trojanGoSni" validate:"omitempty,min=0,max=64"`
	TrojanGoMuxEnable       *uint   `json:"trojanGoMuxEnable" form:"trojanGoMuxEnable" validate:"required,oneof=0 1"`
	TrojanGoWebsocketEnable *uint   `json:"trojanGoWebsocketEnable" form:"trojanGoWebsocketEnable" validate:"required,oneof=0 1"`
	TrojanGoWebsocketPath   *string `json:"trojanGoWebsocketPath" form:"trojanGoWebsocketPath" validate:"omitempty,min=0,max=64"`
	TrojanGoWebsocketHost   *string `json:"trojanGoWebsocketHost" form:"trojanGoWebsocketHost" validate:"omitempty,min=0,max=64"`
	TrojanGoSsEnable        *uint   `json:"trojanGoSsEnable" form:"trojanGoSsEnable" validate:"required,oneof=0 1"`
	TrojanGoSsMethod        *string `json:"trojanGoSsMethod" form:"trojanGoSsMethod" validate:"omitempty,min=0,max=32"`
	TrojanGoSsPassword      *string `json:"trojanGoSsPassword" form:"trojanGoSsPassword" validate:"omitempty,min=0,max=64"`

	HysteriaProtocol *string `json:"hysteriaProtocol" form:"hysteriaProtocol" validate:"omitempty,min=0,max=32"`
	HysteriaUpMbps   *int    `json:"hysteriaUpMbps" form:"hysteriaUpMbps" validate:"required,gt=0,lte=9999999999"`
	HysteriaDownMbps *int    `json:"hysteriaDownMbps" form:"hysteriaDownMbps" validate:"required,gt=0,lte=9999999999"`
}

type NodeUpdateDto struct {
	RequiredIdDto
	NodeServerId *uint   `json:"nodeServerId" form:"nodeServerId" validate:"required,gt=0"`
	NodeSubId    *uint   `json:"nodeSubId" form:"nodeSubId" validate:"required,gte=0"`
	NodeTypeId   *uint   `json:"nodeTypeId" form:"nodeTypeId" validate:"required,gt=0"`
	Name         *string `json:"name" form:"name" validate:"required,min=2,max=20"`
	Domain       *string `json:"domain" form:"domain" validate:"required,ip|fqdn,min=4,max=64"`
	Port         *uint   `json:"port" form:"port" validate:"required,validatePort"`

	XrayFlow           *string `json:"xrayFlow" form:"xrayFlow" validate:"omitempty,min=0,max=32"`
	XraySSMethod       *string `json:"xraySSMethod" form:"xraySSMethod" validate:"omitempty,min=0,max=32"`
	XrayProtocol       *string `json:"xrayProtocol" form:"xrayProtocol" validate:"omitempty,min=0,max=32"`
	XraySettings       *string `json:"xraySettings" form:"xraySettings" validate:"omitempty,min=0,max=256"`
	XrayStreamSettings *string `json:"xrayStreamSettings" form:"xrayStreamSettings" validate:"omitempty,min=0,max=256"`
	XrayTag            *string `json:"xrayTag" form:"xrayTag" validate:"omitempty,min=0,max=64"`
	XraySniffing       *string `json:"xraySniffing" form:"xraySniffing" validate:"omitempty,min=0,max=256"`
	XrayAllocate       *string `json:"xrayAllocate" form:"xrayAllocate" validate:"omitempty,min=0,max=256"`

	TrojanGoSni             *string `json:"trojanGoSni" form:"trojanGoSni" validate:"omitempty,min=0,max=64"`
	TrojanGoMuxEnable       *uint   `json:"trojanGoMuxEnable" form:"trojanGoMuxEnable" validate:"required,oneof=0 1"`
	TrojanGoWebsocketEnable *uint   `json:"trojanGoWebsocketEnable" form:"trojanGoWebsocketEnable" validate:"required,oneof=0 1"`
	TrojanGoWebsocketPath   *string `json:"trojanGoWebsocketPath" form:"trojanGoWebsocketPath" validate:"omitempty,min=0,max=64"`
	TrojanGoWebsocketHost   *string `json:"trojanGoWebsocketHost" form:"trojanGoWebsocketHost" validate:"omitempty,min=0,max=64"`
	TrojanGoSsEnable        *uint   `json:"trojanGoSsEnable" form:"trojanGoSsEnable" validate:"required,oneof=0 1"`
	TrojanGoSsMethod        *string `json:"trojanGoSsMethod" form:"trojanGoSsMethod" validate:"omitempty,min=0,max=32"`
	TrojanGoSsPassword      *string `json:"trojanGoSsPassword" form:"trojanGoSsPassword" validate:"omitempty,min=0,max=64"`

	HysteriaProtocol *string `json:"hysteriaProtocol" form:"hysteriaProtocol" validate:"omitempty,min=0,max=32"`
	HysteriaUpMbps   *int    `json:"hysteriaUpMbps" form:"hysteriaUpMbps" validate:"required,gt=0,lte=9999999999"`
	HysteriaDownMbps *int    `json:"hysteriaDownMbps" form:"hysteriaDownMbps" validate:"required,gt=0,lte=9999999999"`
}
