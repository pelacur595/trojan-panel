package dto

type SendEmailDto struct {
	FromEmailName string   `json:"fromEmailName" form:"fromEmailName" validate:"require,min=1,max=32"`
	ToEmails      []string `json:"toEmails" form:"toEmails" validate:"require,validateEmail"`
	Subject       string   `json:"subject" form:"subject" validate:"require,min=1,max=64"`
	Content       string   `json:"content" form:"content" validate:"require,min=1,max=255"`
}
