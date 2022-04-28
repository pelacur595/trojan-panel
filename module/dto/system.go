package dto

type SystemUpdateDto struct {
	RequiredIdDto
	OpenRegister       *uint   `json:"openRegister" form:"openRegister" validate:"required,oneof=0 1"` // 是否开放注册
	RegisterQuota      *int    `json:"registerQuota" validate:"-"`                                     // 注册用户默认配额 单位/byte
	RegisterExpireDays *uint   `json:"registerExpireDays" validate:"omitempty,gte=0"`                  // 注册用户过期天数 单位/天
	EmailHost          *string `json:"emailHost" form:"emailHost" validate:"required,hostname|fqdn,min=3,max=64"`
	EmailPort          *uint   `json:"emailPort" form:"emailPort" validate:"required,gt=0,lte=65535"`
	EmailUsername      *string `json:"emailUsername" form:"emailUsername" validate:"required,min=3,max=32"`
	EmailPassword      *string `json:"emailPassword" form:"emailPassword" validate:"required,min=3,max=32"`
}
