package dto

import "time"

type RegisterInput struct {
	// User Info
	Name        string    `json:"name" validate:"required"`
	DateOfBirth time.Time `json:"date_of_birth" validate:"required"` 
	Gender      string    `json:"gender" validate:"required,oneof=male female"`
		BloodGroup  string    `json:"blood_group" validate:"required"` // e.g. "O+"
	PhoneNumber string    `json:"phone_number" validate:"required"`
	Email       string    `json:"email" validate:"omitempty,email"` 
	UserName    string    `json:"username" validate:"required"`
	
	// Address Info 
	Area       string `json:"area" validate:"required"`
	City       string `json:"city" validate:"required"`
	State      string `json:"state" validate:"required"`
	Country    string `json:"country" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required"`
	
	// Auth
	Password string `json:"password" validate:"required,min=4,max=12"`
}

type LoginInput struct {
	UserName string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}