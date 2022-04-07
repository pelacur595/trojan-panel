package dto

type SystemUpdateDto struct {
	RequiredIdDto
	OpenRegister       *uint `json:"openRegister" form:"openRegister" validate:"required,oneof=0 1"` // 是否开放注册
	RegisterQuota      *int  `json:"registerQuota" validate:"-"`                                     // 注册用户默认配额 单位/byte
	RegisterExpireDays *uint `json:"registerExpireDays" validate:"omitempty,gte=0"`                  // 注册用户过期天数 单位/天
}
