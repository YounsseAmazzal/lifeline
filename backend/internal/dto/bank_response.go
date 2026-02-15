package dto

import "time"

type BankResponse struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	City        string          `json:"city"`
	PhoneNumber string          `json:"phoneNumber"`
	Email       string          `json:"email"`
	Website     string          `json:"website"`
	PhotoURL    string          `json:"photoUrl"`
	LastUpdated time.Time       `json:"lastUpdated"`
	Latitude    float64         `json:"latitude"`  
	Longitude   float64         `json:"longitude"` 
	
	BloodGroups []BloodGroupDto `json:"bloodGroups"`
}

type BloodGroupDto struct {
	ID       uint   `json:"id"`
	Group    string `json:"group"`
	Quantity int    `json:"quantity"`
}