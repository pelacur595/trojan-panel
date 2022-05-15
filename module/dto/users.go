package dto

type UsersPageDto struct {
	BaseDto
	UsersDto
}

type UsersDto struct {
	Username *string `json:"username" form:"username" validate:"omitempty,min=0,max=20"`
}

type UserUpdateProfileDto struct {
	Username *string `json:"username" form:"username" validate:"required,min=6,max=20,validateStr"`
	NewPass  *string `json:"newPass" form:"newPass" validate:"required,min=6,max=20,validateStr"`
	OldPass  *string `json:"oldPass" form:"oldPass" validate:"required,min=6,max=20,validateStr"`
	Email    *string `json:"email" form:"email" validate:"omitempty,validateEmail"`
}

type UserCreateDto struct {
	Quota      *int    `json:"quota" form:"quota" validate:"required,gte=-1,lte=1024000"`
	Username   *string `json:"username" form:"username" validate:"required,min=6,max=20,validateStr"`
	Pass       *string `json:"pass" form:"pass" validate:"required,min=6,max=20,validateStr"`
	RoleId     *uint   `json:"roleId" form:"roleId" validate:"required,oneof=2 3"`
	Deleted    *uint   `json:"deleted" form:"deleted" validate:"required,oneof=0 1"`
	ExpireTime *uint   `json:"expireTime" form:"expireTime" validate:"required"`
	Email      *string `json:"email" form:"email" validate:"omitempty,validateEmail"`
	//IpLimit            *uint   `json:"ipLimit" form:"ipLimit" validate:"required,gt=0"`
	//UploadSpeedLimit   *uint   `json:"uploadSpeedLimit" form:"uploadSpeedLimit" validate:"required"`
	//DownloadSpeedLimit *uint   `json:"downloadSpeedLimit" form:"downloadSpeedLimit" validate:"required"`
}

type UserUpdateDto struct {
	RequiredIdDto
	Quota      *int    `json:"quota" form:"quota" validate:"required,gte=-1,lte=1024000"`
	Username   *string `json:"username" form:"username" validate:"omitempty,min=0,max=20,validateStr"`
	Pass       *string `json:"pass" form:"pass" validate:"omitempty,min=6,max=20,validateStr"`
	RoleId     *uint   `json:"roleId" form:"roleId" validate:"required,oneof=1 2 3"`
	Deleted    *uint   `json:"deleted" form:"deleted" validate:"required,oneof=0 1"`
	ExpireTime *uint   `json:"expireTime" form:"expireTime" validate:"required,gte=0"`
	Email      *string `json:"email" form:"email" validate:"omitempty,validateEmail"`
	//IpLimit            *uint   `json:"ipLimit" form:"ipLimit" validate:"required,gt=0"`
	//UploadSpeedLimit   *uint   `json:"uploadSpeedLimit" form:"uploadSpeedLimit" validate:"required"`
	//DownloadSpeedLimit *uint   `json:"downloadSpeedLimit" form:"downloadSpeedLimit" validate:"required"`
}

type UserLoginDto struct {
	Username *string `json:"username" form:"username" validate:"required,min=6,max=20,validateStr"`
	Pass     *string `json:"pass" form:"pass" validate:"required,min=6,max=20,validateStr"`
}

type UserRegisterDto struct {
	Username *string `json:"username" form:"username" validate:"required,min=6,max=20,validateStr,excludes=admin"`
	Pass     *string `json:"pass" form:"pass" validate:"required,min=6,max=20,validateStr"`
}

type HysteriaAutoDto struct {
	Payload *string `json:"payload" validate:"required"`
}
