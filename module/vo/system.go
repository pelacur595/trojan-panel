package vo

type SystemVo struct {
	Id                 uint   `json:"id" redis:"id"`
	OpenRegister       uint   `json:"openRegister" redis:"openRegister"`
	RegisterQuota      int    `json:"registerQuota" redis:"registerQuota"`
	RegisterExpireDays uint   `json:"registerExpireDays" redis:"registerExpireDays"`
	EmailHost          string `json:"emailHost" redis:"emailHost"`
	EmailPort          uint   `json:"emailPort" redis:"emailPort"`
	EmailUsername      string `json:"emailUsername" redis:"emailUsername"`
	EmailPassword      string `json:"emailPassword" redis:"emailPassword"`
}

type SettingVo struct {
	OpenRegister uint `json:"openRegister"`
}
