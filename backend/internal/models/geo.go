package models

type Region struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Region    string `json:"region"` // e.g. "Tanger-Tetouan-Al Hoceima"
	Cities    []City `json:"cities" gorm:"foreignKey:RegionID"`
}

type City struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Ville    string `json:"ville"`   // e.g. "Tanger"
	RegionID uint   `json:"region_id"`
}