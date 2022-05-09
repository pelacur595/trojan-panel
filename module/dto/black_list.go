package dto

type BlackListDto struct {
	Ip *string `json:"ip" form:"ip" validate:"omitempty,min=0,max=64"`
}

type BlackListPageDto struct {
	BaseDto
	BlackListDto
}

type BlackListCreateDto struct {
	Ips []string `json:"ip" form:"ip" validate:"required,hostname|fqdn,min=3,max=64"`
}

type BlackListDeleteDto struct {
	Ip *string `json:"ip" form:"ip" validate:"required,hostname|fqdn,min=3,max=64"`
}
