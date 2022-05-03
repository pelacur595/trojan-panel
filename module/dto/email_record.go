package dto

type EmailRecordPageDto struct {
	BaseDto
	EmailRecordDto
}

type EmailRecordDto struct {
	ToEmail *string `json:"toEmail" form:"toEmail" validate:"omitempty,validateEmail"`
	State   *int    `json:"state" form:"state" validate:"required,oneof=-1 0 1"`
}

type EmailRecordUpdateDto struct {
	RequiredIdDto
	State *int `json:"state" form:"state" validate:"required,oneof=-1 0 1"`
}
