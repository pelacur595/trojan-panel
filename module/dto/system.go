package dto

type SystemUpdateDto struct {
	RequiredIdDto
	RegisterEnable              *uint `json:"registerEnable" form:"registerEnable" validate:"required,oneof=0 1"`
	RegisterQuota               *int  `json:"registerQuota" validate:"required,gte=0,lte=1024000"`
	RegisterExpireDays          *uint `json:"registerExpireDays" validate:"required,gte=0,lte=365"`
	ResetDownloadAndUploadMonth *uint `json:"resetDownloadAndUploadMonth" validate:"required,oneof=0 1"`
	TrafficRankEnable           *uint `json:"trafficRankEnable" validate:"required,oneof=0 1"`

	ExpireWarnEnable *uint   `json:"expireWarnEnable" redis:"expireWarnEnable" validate:"required,oneof=0 1"`
	ExpireWarnDay    *uint   `json:"expireWarnDay" redis:"expireWarnDay" validate:"required,gte=0,lte=365"`
	EmailEnable      *uint   `json:"emailEnable" redis:"emailEnable" validate:"required,oneof=0 1"`
	EmailHost        *string `json:"emailHost" form:"emailHost" validate:"omitempty,min=0,max=64"`
	EmailPort        *uint   `json:"emailPort" form:"emailPort" validate:"required,validatePort"`
	EmailUsername    *string `json:"emailUsername" form:"emailUsername" validate:"omitempty,min=0,max=32"`
	EmailPassword    *string `json:"emailPassword" form:"emailPassword" validate:"omitempty,min=0,max=32"`

	SystemName *string `json:"systemName" form:"systemName" validate:"omitempty,min=0,max=32"`
	ClashRule  *string `json:"clashRule" form:"clashRule" validate:"omitempty,min=0,max=5120"`
}
