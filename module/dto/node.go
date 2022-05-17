package dto

type NodePageDto struct {
	NodeDto
	BaseDto
}

type NodeDto struct {
	Name *string `json:"name" form:"name" validate:"omitempty,min=0,max=20"`
}

type NodeQRCodeDto struct {
	Name             *string `json:"name" form:"name" validate:"required,min=2,max=20"`
	Ip               *string `json:"ip" form:"ip" validate:"required,ip|fqdn,min=4,max=64"`
	Port             *uint   `json:"port" form:"port" validate:"required,validatePort"`
	Type             *uint   `json:"type" form:"type" validate:"required,gt=0"`
	WebsocketEnable  *uint   `json:"websocketEnable" validate:"required,oneof=0 1"`
	WebsocketPath    *string `json:"websocketPath" validate:"omitempty,min=0,max=32"`
	SsEnable         *uint   `json:"ssEnable" validate:"required,oneof=0 1"`
	SsMethod         *string `json:"ssMethod" validate:"omitempty,min=0,max=16"`
	SsPassword       *string `json:"ssPassword" validate:"omitempty,min=0,max=32"`
	HysteriaProtocol *string `json:"hysteriaProtocol" validate:"required,min=0,max=16"`
	HysteriaUpMbps   *int    `json:"hysteriaUpMbps" validate:"required,gt=0,lte=999"`
	HysteriaDownMbps *int    `json:"hysteriaDownMbps" validate:"required,gt=0,lte=999"`
}

type NodeCreateDto struct {
	Name             *string `json:"name" form:"name" validate:"required,min=2,max=20"`
	Ip               *string `json:"ip" form:"ip" validate:"required,ip|fqdn,min=4,max=64"`
	Port             *uint   `json:"port" form:"port" validate:"required,validatePort"`
	Type             *uint   `json:"type" form:"type" validate:"required,gt=0"`
	WebsocketEnable  *uint   `json:"websocketEnable" validate:"required,oneof=0 1"`
	WebsocketPath    *string `json:"websocketPath" validate:"omitempty,min=0,max=32"`
	SsEnable         *uint   `json:"ssEnable" validate:"required,oneof=0 1"`
	SsMethod         *string `json:"ssMethod" validate:"omitempty,min=0,max=16"`
	SsPassword       *string `json:"ssPassword" validate:"omitempty,min=0,max=32"`
	HysteriaProtocol *string `json:"hysteriaProtocol" validate:"required,min=0,max=16"`
	HysteriaUpMbps   *int    `json:"hysteriaUpMbps" validate:"required,gt=0,lte=999"`
	HysteriaDownMbps *int    `json:"hysteriaDownMbps" validate:"required,gt=0,lte=999"`
}

type NodeUpdateDto struct {
	RequiredIdDto
	Name             *string `json:"name" form:"name" validate:"required,min=2,max=20"`
	Ip               *string `json:"ip" form:"ip" validate:"required,ip|fqdn,min=4,max=64"`
	Port             *uint   `json:"port" form:"port" validate:"required,validatePort"`
	Type             *uint   `json:"type" form:"type" validate:"required,gt=0"`
	WebsocketEnable  *uint   `json:"websocketEnable" validate:"required,oneof=0 1"`
	WebsocketPath    *string `json:"websocketPath" validate:"omitempty,min=0,max=32"`
	SsEnable         *uint   `json:"ssEnable" validate:"required,oneof=0 1"`
	SsMethod         *string `json:"ssMethod" validate:"omitempty,min=0,max=16"`
	SsPassword       *string `json:"ssPassword" validate:"omitempty,min=0,max=32"`
	HysteriaProtocol *string `json:"hysteriaProtocol" validate:"required,min=0,max=16"`
	HysteriaUpMbps   *int    `json:"hysteriaUpMbps" validate:"required,gt=0,lte=999"`
	HysteriaDownMbps *int    `json:"hysteriaDownMbps" validate:"required,gt=0,lte=999"`
}
