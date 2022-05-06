package dto

type SendEmailDto struct {
	EmailUsername string   `json:"emailUsername" form:"emailUsername" validate:"require,min=0,max=32"`
	FromEmailName string   `json:"fromEmailName" form:"fromEmailName" validate:"require,min=0,max=32"`
	ToEmails      []string `json:"toEmails" form:"toEmails" validate:"require,validateEmail"`
	Subject       string   `json:"subject" form:"subject" validate:"require,min=0,max=64"`
	Content       string   `json:"content" form:"content" validate:"require,min=0,max=255"`
}
