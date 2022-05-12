package dto

type BaseDto struct {
	PageNum   *uint `json:"pageNum" form:"pageNum" validate:"required,validatePositiveInt"`
	PageSize  *uint `json:"pageSize" form:"pageSize" validate:"required,validatePositiveInt"`
	StartTime *uint `json:"startTime" form:"startTime" validate:"omitempty,validatePositiveInt"`
	EndTime   *uint `json:"endTime" form:"endTime" validate:"omitempty,validatePositiveInt"`
}

type RequiredIdDto struct {
	Id *uint `json:"id" form:"id" validate:"required,validatePositiveInt"`
}
