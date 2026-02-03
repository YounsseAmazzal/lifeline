package models

import (
	"time"
)

type Sponsor struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	
	Name        string    `json:"name" gorm:"unique;not null"` // e.g. "Banque Populaire"
	Description string    `json:"description"`                 // e.g. "Official Partner 2024"
	LogoURL     string    `json:"logo_url"`                    
	Website     string    `json:"website"`                     

	TotalPaid      float64   `json:"total_paid"`      
	CampaignBudget float64   `json:"campaign_budget"` 
	IsActive       bool      `json:"is_active" gorm:"default:true"` 

	
	ViewsCount     int64     `json:"views_count" gorm:"default:0"`  
	ClicksCount    int64     `json:"clicks_count" gorm:"default:0"` 

	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}