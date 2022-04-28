package vo

type SystemVo struct {
	Id                 uint   `json:"id"`
	OpenRegister       uint   `json:"openRegister"`
	RegisterQuota      int    `json:"registerQuota"`
	RegisterExpireDays uint   `json:"registerExpireDays"`
	EmailHost          string `json:"emailHost"`
	EmailPort          uint   `json:"emailPort"`
	EmailUsername      string `json:"emailUsername"`
	EmailPassword      string `json:"emailPassword"`
}

type SettingVo struct {
	OpenRegister uint `json:"openRegister"`
}
