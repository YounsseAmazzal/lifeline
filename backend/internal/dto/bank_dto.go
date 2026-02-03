package dto

// BankRegisterInput:
type BankRegisterInput struct {
	Name string `json:"name" validate:"required"`

	// Address Info
	Area       string `json:"area" validate:"required"`
	City       string `json:"city" validate:"required"`
	State      string `json:"state" validate:"required"`
	Country    string `json:"country" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required"`

	// Contact Info
	PhoneNumber string `json:"phone_number" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Website     string `json:"website"` 

	// Admin Logic
	BankAdmin string `json:"bank_admin" validate:"required"` 
}