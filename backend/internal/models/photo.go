package models

type Photo struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	PublicID string `json:"public_id"`
	URL      string `json:"url"`

	// Relationships (Nullable)
	UserID *uint `json:"user_id"`
	BankID *uint `json:"bank_id"`
}