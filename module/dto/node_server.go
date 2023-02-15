package dto

type NodeServerPageDto struct {
	NodeServerDto
	BaseDto
}

type NodeServerDto struct {
	Name *string `json:"name" form:"name" validate:"omitempty,min=0,max=20"`
	Ip   *string `json:"ip" form:"ip" validate:"omitempty,min=0,max=64"`
}

type NodeServerCreateDto struct {
	Name     *string `json:"name" form:"name" validate:"required,min=2,max=20"`
	Ip       *string `json:"ip" form:"ip" validate:"required,ip|fqdn,min=4,max=64"`
	GrpcPort *uint   `json:"grpcPort" form:"grpcPort" validate:"required,validatePort"`
}

type NodeServerUpdateDto struct {
	RequiredIdDto
	Name     *string `json:"name" form:"name" validate:"required,min=2,max=20"`
	Ip       *string `json:"ip" form:"ip" validate:"required,ip|fqdn,min=4,max=64"`
	GrpcPort *uint   `json:"grpcPort" form:"grpcPort" validate:"required,validatePort"`
}
