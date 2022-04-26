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

type SystemEmailVo struct {
	Id            uint   `json:"id"`
	EmailHost     string `json:"emailHost"`
	EmailPort     int    `json:"emailPort"`
	EmailUsername string `json:"emailUsername"`
	EmailPassword string `json:"emailPassword"`
}
