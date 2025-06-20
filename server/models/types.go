package models

type BillingStatus string

const (
	Billed    BillingStatus = "Billed"
	NotBilled BillingStatus = "NotBilled"
)

type SignupInput struct {
	Nickname string `json:"nickname" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Avatar   string `json:"avatar" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}
