package dto

type SystemUpdateDto struct {
	RequiredIdDto
	OpenRegister       *uint   `json:"openRegister" form:"openRegister" validate:"required,oneof=0 1"` // 是否开放注册
	RegisterQuota      *int    `json:"registerQuota" validate:"required,gte=0,lte=1024000"`            // 注册用户默认配额 单位/MB
	RegisterExpireDays *uint   `json:"registerExpireDays" validate:"required,gte=0,lte=365"`           // 注册用户过期天数 单位/天
	ExpireWarnEnable   *uint   `json:"expireWarnEnable" redis:"expireWarnEnable" validate:"required,oneof=0 1"`
	ExpireWarnDay      *uint   `json:"expireWarnDay" redis:"expireWarnDay" validate:"required,gte=0,lte=365"`
	EmailEnable        *uint   `json:"emailEnable" redis:"emailEnable" validate:"required,oneof=0 1"`
	EmailHost          *string `json:"emailHost" form:"emailHost" validate:"omitempty,ip|fqdn,min=0,max=64"`
	EmailPort          *uint   `json:"emailPort" form:"emailPort" validate:"required,validatePort"`
	EmailUsername      *string `json:"emailUsername" form:"emailUsername" validate:"omitempty,min=0,max=32"`
	EmailPassword      *string `json:"emailPassword" form:"emailPassword" validate:"omitempty,min=0,max=32"`
}
