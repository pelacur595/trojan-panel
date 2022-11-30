package bo

type SystemRegisterConfigBo struct {
	OpenRegister       uint `json:"openRegister"`       // 是否开放注册
	RegisterQuota      int  `json:"registerQuota"`      // 注册用户默认配额 单位/MB
	RegisterExpireDays uint `json:"registerExpireDays"` // 注册用户过期天数 单位/天
}

type SystemEmailConfigBo struct {
	ExpireWarnEnable uint   `json:"expireWarnEnable"`
	ExpireWarnDay    uint   `json:"expireWarnDay"`
	EmailEnable      uint   `json:"emailEnable"`
	EmailHost        string `json:"emailHost"`
	EmailPort        uint   `json:"emailPort"`
	EmailUsername    string `json:"emailUsername"`
	EmailPassword    string `json:"emailPassword"`
}
