package models

type BloodGroup struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Group string `json:"group"`    // e.g. "O+"
	Quantity int    `json:"quantity"` 
	BankID uint `json:"bank_id"`
}