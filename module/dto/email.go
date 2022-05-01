package dto

type SendEmailDto struct {
	FromEmailName string   `json:"fromEmailName" form:"fromEmailName" validate:"omitempty,min=0,max=32"`
	ToEmails      []string `json:"toEmails" form:"toEmails" validate:"omitempty,validateEmail"`
	Subject       string   `json:"subject" form:"subject" validate:"omitempty,min=0,max=64"`
	Content       string   `json:"content" form:"content" validate:"omitempty,min=0,max=255"`
}
