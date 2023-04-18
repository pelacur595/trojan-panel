package dto

type AccountPageDto struct {
	BaseDto
	AccountDto
}

type AccountDto struct {
	Username *string `json:"username" form:"username" validate:"omitempty,min=0,max=20"`
}

type AccountUpdateProfileDto struct {
	Username *string `json:"username" form:"username" validate:"required,min=6,max=20,validateStr"`
	NewPass  *string `json:"newPass" form:"newPass" validate:"required,min=6,max=20,validateStr"`
	OldPass  *string `json:"oldPass" form:"oldPass" validate:"required,min=6,max=20,validateStr"`
	Email    *string `json:"email" form:"email" validate:"omitempty,validateEmail"`
}

type AccountCreateDto struct {
	Username   *string `json:"username" form:"username" validate:"required,min=6,max=20,validateStr"`
	Pass       *string `json:"pass" form:"pass" validate:"required,min=6,max=20,validateStr"`
	RoleId     *uint   `json:"roleId" form:"roleId" validate:"required,oneof=2 3"`
	Email      *string `json:"email" form:"email" validate:"omitempty,validateEmail"`
	ExpireTime *uint   `json:"expireTime" form:"expireTime" validate:"required"`
	Deleted    *uint   `json:"deleted" form:"deleted" validate:"required,oneof=0 1"`
	Quota      *int    `json:"quota" form:"quota" validate:"required,gte=-1,lte=1024000"`
	//IpLimit            *uint   `json:"ipLimit" form:"ipLimit" validate:"required,gt=0"`
	//UploadSpeedLimit   *uint   `json:"uploadSpeedLimit" form:"uploadSpeedLimit" validate:"required"`
	//DownloadSpeedLimit *uint   `json:"downloadSpeedLimit" form:"downloadSpeedLimit" validate:"required"`
}

type AccountUpdateDto struct {
	RequiredIdDto
	Quota      *int    `json:"quota" form:"quota" validate:"required,gte=-1,lte=1024000"`
	Username   *string `json:"username" form:"username" validate:"required,min=0,max=20,validateStr"`
	Pass       *string `json:"pass" form:"pass" validate:"omitempty,min=6,max=20,validateStr"`
	RoleId     *uint   `json:"roleId" form:"roleId" validate:"required,oneof=1 2 3"`
	Deleted    *uint   `json:"deleted" form:"deleted" validate:"required,oneof=0 1"`
	ExpireTime *uint   `json:"expireTime" form:"expireTime" validate:"required,gte=0"`
	Email      *string `json:"email" form:"email" validate:"omitempty,validateEmail"`
	//IpLimit            *uint   `json:"ipLimit" form:"ipLimit" validate:"required,gt=0"`
	//UploadSpeedLimit   *uint   `json:"uploadSpeedLimit" form:"uploadSpeedLimit" validate:"required"`
	//DownloadSpeedLimit *uint   `json:"downloadSpeedLimit" form:"downloadSpeedLimit" validate:"required"`
}

type AccountLoginDto struct {
	Username *string `json:"username" form:"username" validate:"required,min=6,max=20,validateStr"`
	Pass     *string `json:"pass" form:"pass" validate:"required,min=6,max=20,validateStr"`
	CaptureDto
}

type AccountRegisterDto struct {
	Username *string `json:"username" form:"username" validate:"required,min=6,max=20,validateStr,excludes=admin"`
	Pass     *string `json:"pass" form:"pass" validate:"required,min=6,max=20,validateStr"`
	CaptureDto
}

type HysteriaAutoDto struct {
	Payload *string `json:"payload" form:"payload" validate:"required"`
}

type CaptureDto struct {
	CaptchaId   *string `json:"captchaId" form:"captchaId" validate:"required"`
	CaptchaCode *string `json:"captchaCode" form:"captchaCode" validate:"required"`
}
