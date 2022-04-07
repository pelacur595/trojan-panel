package vo

type SystemVo struct {
	Id                 uint `json:"id"`
	OpenRegister       uint `json:"openRegister"`
	RegisterQuota      int  `json:"registerQuota"`
	RegisterExpireDays uint `json:"registerExpireDays"`
}

type SettingVo struct {
	OpenRegister uint `json:"openRegister"`
}
