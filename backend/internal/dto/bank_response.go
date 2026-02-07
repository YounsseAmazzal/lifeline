package dto

import "time"

type BankResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	
	// Address Flattened
	Area       string `json:"area"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`

	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Website     string    `json:"website"`
	LastUpdated time.Time `json:"last_updated"`
	PhotoURL    string    `json:"photo_url"`
	Latitude    float64         `json:"latitude"`  
	Longitude   float64         `json:"longitude"` 
	BloodGroups []BloodGroupDto `json:"blood_groups"`
}

type BloodGroupDto struct {
	ID       uint   `json:"id"`
	Group    string `json:"group"`    
	Quantity int    `json:"quantity"` 
}