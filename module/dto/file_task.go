package dto

type FileTaskPageDto struct {
	BaseDto
	FileTaskDto
}

type FileTaskDto struct {
	Type *uint `json:"type" form:"type" validate:"omitempty,oneof=0 1 2"`
}
