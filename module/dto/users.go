package dto

type UsersPageDto struct {
	BaseDto
	UsersDto
}

type UsersDto struct {
	Username *string `json:"username" form:"username" validate:"omitempty,min=0,max=20"`
}

type UserUpdatePassDto struct {
	RequiredIdDto
	NewPass *string `json:"newPass" form:"newPass" validate:"required,min=6,max=20,validateStr"`
	OldPass *string `json:"oldPass" form:"oldPass" validate:"required,min=6,max=20,validateStr"`
}

type UserCreateDto struct {
	Quota      *int    `json:"quota" form:"quota" validate:"-"`
	Username   *string `json:"username" form:"username" validate:"required,min=6,max=20,validateStr"`
	Pass       *string `json:"pass" form:"pass" validate:"required,min=6,max=20,validateStr"`
	RoleId     *uint   `json:"roleId" form:"roleId" validate:"required,oneof=2 3"`
	Deleted    *uint   `json:"deleted" form:"deleted" validate:"required,oneof=0 1"`
	ExpireTime *uint   `json:"expireTime" form:"expireTime" validate:"required,gte=0"`
}

type UserUpdateDto struct {
	RequiredIdDto
	Quota      *int    `json:"quota" form:"quota" validate:"-"`
	Username   *string `json:"username" form:"username" validate:"omitempty,min=6,max=20,validateStr"`
	Pass       *string `json:"pass" form:"pass" validate:"omitempty,min=6,max=20,validateStr"`
	RoleId     *uint   `json:"roleId" form:"roleId" validate:"required,oneof=1 2 3"`
	Deleted    *uint   `json:"deleted" form:"deleted" validate:"required,oneof=0 1"`
	ExpireTime *uint   `json:"expireTime" form:"expireTime" validate:"required,gte=0"`
}

type UserLoginDto struct {
	Username *string `json:"username" form:"username" validate:"required,min=6,max=20,validateStr"`
	Pass     *string `json:"pass" form:"pass" validate:"required,min=6,max=20,validateStr"`
}

type UserRegisterDto struct {
	Username *string `json:"username" form:"username" validate:"required,min=6,max=20,validateStr,excludes=admin"`
	Pass     *string `json:"pass" form:"pass" validate:"required,min=6,max=20,validateStr"`
}
