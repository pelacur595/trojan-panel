package dto

type NodePageDto struct {
	NodeDto
	BaseDto
}

type NodeDto struct {
	Name *string `json:"name" form:"name" validate:"omitempty,min=0,max=20"`
}

type NodeQRCodeDto struct {
	Name            *string `json:"name" form:"name" validate:"omitempty,min=3,max=20"`
	Ip              *string `json:"ip" form:"ip" validate:"required,ip|fqdn,min=3,max=64"`
	Port            *uint   `json:"port" form:"port" validate:"required,gt=0,lte=65535,validatePositiveInt"`
	Type            *uint   `json:"type" form:"type" validate:"required,validatePositiveInt"`
	WebsocketEnable *uint   `json:"websocketEnable" validate:"omitempty,oneof=0 1"`
	WebsocketPath   *string `json:"websocketPath" validate:"omitempty,min=0,max=32"`
	SsEnable        *uint   `json:"ssEnable" validate:"omitempty,oneof=0 1"`
	SsMethod        *string `json:"ssMethod" validate:"omitempty,min=0,max=16"`
	SsPassword      *string `json:"ssPassword" validate:"omitempty,min=0,max=32"`
}

type NodeCreateDto struct {
	Name            *string `json:"name" form:"name" validate:"required,min=3,max=20"`
	Ip              *string `json:"ip" form:"ip" validate:"required,ip|fqdn,min=3,max=64"`
	Port            *uint   `json:"port" form:"port" validate:"required,gt=0,lte=65535,validatePositiveInt"`
	Type            *uint   `json:"type" form:"type" validate:"required,validatePositiveInt"`
	WebsocketEnable *uint   `json:"websocketEnable" validate:"omitempty,oneof=0 1"`
	WebsocketPath   *string `json:"websocketPath" validate:"omitempty,min=0,max=32"`
	SsEnable        *uint   `json:"ssEnable" validate:"omitempty,oneof=0 1"`
	SsMethod        *string `json:"ssMethod" validate:"omitempty,min=0,max=16"`
	SsPassword      *string `json:"ssPassword" validate:"omitempty,min=0,max=32"`
}

type NodeUpdateDto struct {
	RequiredIdDto
	Name            *string `json:"name" form:"name" validate:"omitempty,min=3,max=20"`
	Ip              *string `json:"ip" form:"ip" validate:"omitempty,ip|fqdn,min=3,max=64"`
	Port            *uint   `json:"port" form:"port" validate:"omitempty,gt=0,lte=65535,validatePositiveInt"`
	Type            *uint   `json:"type" form:"type" validate:"omitempty,validatePositiveInt"`
	WebsocketEnable *uint   `json:"websocketEnable" validate:"omitempty,oneof=0 1"`
	WebsocketPath   *string `json:"websocketPath" validate:"omitempty,min=0,max=32"`
	SsEnable        *uint   `json:"ssEnable" validate:"omitempty,oneof=0 1"`
	SsMethod        *string `json:"ssMethod" validate:"omitempty,min=0,max=16"`
	SsPassword      *string `json:"ssPassword" validate:"omitempty,min=0,max=32"`
}
