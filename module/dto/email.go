package dto

type SendEmailDto struct {
	FromEmailName string `json:"fromEmailName" form:"fromEmailName" validate:"omitempty,min=0,max=32"`
	ToEmail       string `json:"toEmail" form:"toEmail" validate:"omitempty,validateEmail"`
	Subject       string `json:"subject" form:"subject" validate:"omitempty,min=0,max=64"`
	Content       string `json:"content" form:"content" validate:"omitempty,min=0,max=255"`
}
