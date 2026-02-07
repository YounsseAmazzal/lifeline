package dto

type RegisterInput struct {
	Name        string `json:"name" validate:"required"`
	UserName    string `json:"userName" validate:"required"`
	Email       string `json:"email" validate:"omitempty,email"`
	Password    string `json:"password" validate:"required,min=6"` 
	PhoneNumber string `json:"phoneNumber" validate:"required"`   
	BloodGroup  string `json:"bloodGroup" validate:"required"`   
	Gender      string `json:"gender" validate:"omitempty"`       

	DateOfBirth string `json:"dateOfBirth" validate:"omitempty"`

	City       string `json:"city" validate:"required"`
	Area       string `json:"area" validate:"omitempty"`
	State      string `json:"state" validate:"omitempty"`
	Country    string `json:"country" validate:"omitempty"`
	PostalCode string `json:"postalCode" validate:"omitempty"`
}

type LoginInput struct {
	UserName string `json:"userName" validate:"required"` 
	Password string `json:"password" validate:"required"`
}
