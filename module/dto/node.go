package dto

type NodePageDto struct {
	NodeDto
	BaseDto
}

type NodeDto struct {
	Name *string `json:"name" form:"name" validate:"omitempty,min=0,max=20"`
}

type NodeQRCodeDto struct {
	Name *string `json:"name" form:"name" validate:"omitempty,min=3,max=20"`
	Ip   *string `json:"ip" form:"ip" validate:"required,hostname|fqdn"`
	Port *uint   `json:"port" form:"port" validate:"required,gt=0,lte=65535"`
	Type *uint   `json:"type" form:"type" validate:"required,validatePositiveInt"`
}

type NodeCreateDto struct {
	Name *string `json:"name" form:"name" validate:"required,min=3,max=20"`
	Ip   *string `json:"ip" form:"ip" validate:"required,hostname|fqdn"`
	Port *uint   `json:"port" form:"port" validate:"required,gt=0,lte=65535"`
	Type *uint   `json:"type" form:"type" validate:"required,validatePositiveInt"`
}

type NodeUpdateDto struct {
	RequiredIdDto
	Name *string `json:"name" form:"name" validate:"omitempty,min=3,max=20"`
	Ip   *string `json:"ip" form:"ip" validate:"omitempty,hostname|fqdn"`
	Port *uint   `json:"port" form:"port" validate:"required,gt=0,lte=65535"`
	Type *uint   `json:"type" form:"type" validate:"required,validatePositiveInt"`
}
