package models

type Moderator struct {
	// Composite Primary Key (UserID + BankID)
	UserID uint `gorm:"primaryKey" json:"user_id"`
	BankID uint `gorm:"primaryKey" json:"bank_id"`

	// Additional Info
	Type string `json:"type"` // e.g. "Admin", "Manager"

	// Associations (Optional for preloading)
	// User User `gorm:"foreignKey:UserID"`
	// Bank Bank `gorm:"foreignKey:BankID"`
}