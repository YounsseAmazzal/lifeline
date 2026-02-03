package dto

import "time"

type UserProfile struct {
	ID          uint      `json:"id"` // Read-only usually
	Name        string    `json:"name" validate:"required"`
	DateOfBirth time.Time `json:"date_of_birth" validate:"required"`
	Gender      string    `json:"gender" validate:"required"`
	BloodGroup  string    `json:"blood_group" validate:"required"`
	PhoneNumber string    `json:"phone_number" validate:"required"`
	Email       string    `json:"email" validate:"omitempty,email"`
	UserName    string    `json:"username" validate:"required"`
	
	// Address Fields
	Area       string `json:"area" validate:"required"`
	City       string `json:"city" validate:"required"`
	State      string `json:"state" validate:"required"`
	Country    string `json:"country" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required"`
	
	// Status & Photo
	Available bool   `json:"available"` 
	PhotoURL  string `json:"photo_url"` 
}