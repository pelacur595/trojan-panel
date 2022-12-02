package vo

type SystemVo struct {
	Id                          uint `json:"id" redis:"id"`
	RegisterEnable              uint `json:"registerEnable" redis:"registerEnable"`
	RegisterQuota               int  `json:"registerQuota" redis:"registerQuota"`
	RegisterExpireDays          uint `json:"registerExpireDays" redis:"registerExpireDays"`
	ResetDownloadAndUploadMonth uint `json:"resetDownloadAndUploadMonth" redis:"resetDownloadAndUploadMonth"`
	TrafficRankEnable           uint `json:"trafficRankEnable" redis:"trafficRankEnable"`

	ExpireWarnEnable uint   `json:"expireWarnEnable" redis:"expireWarnEnable"`
	ExpireWarnDay    uint   `json:"expireWarnDay" redis:"expireWarnDay"`
	EmailEnable      uint   `json:"emailEnable"`
	EmailHost        string `json:"emailHost" redis:"emailHost"`
	EmailPort        uint   `json:"emailPort" redis:"emailPort"`
	EmailUsername    string `json:"emailUsername" redis:"emailUsername"`
	EmailPassword    string `json:"emailPassword" redis:"emailPassword"`
}

type SettingVo struct {
	RegisterEnable     uint `json:"registerEnable"`
	RegisterQuota      int  `json:"registerQuota"`
	RegisterExpireDays uint `json:"registerExpireDays"`
	TrafficRankEnable  uint `json:"trafficRankEnable"`
	EmailEnable        uint `json:"emailEnable"`
}
