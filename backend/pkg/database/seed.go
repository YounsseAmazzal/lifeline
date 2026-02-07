package database

import (
	"encoding/json"
	"fmt"
	"lifeline/internal/models"
	"os" 

	"gorm.io/gorm"
)

// --- Structs for JSON Parsing ---
type JsonRegion struct {
	ID     string `json:"id"`
	Region string `json:"region"`
}

type JsonCity struct {
	ID       string `json:"id"`
	Ville    string `json:"ville"`
	RegionID string `json:"region"`
}

type JsonAddress struct {
	Area       string `json:"Area"`
	City       string `json:"City"`
	State      string `json:"State"`
	PostalCode string `json:"PostalCode"`
}

type JsonBloodGroup struct {
	Group string `json:"Group"`
	Value int    `json:"Value"`
}

type JsonBank struct {
	Name        string           `json:"Name"`
	PhoneNumber string           `json:"PhoneNumber"`
	Email       string           `json:"Email"`
	Website     string           `json:"Website"`
	Address     JsonAddress      `json:"Address"`
	BloodGroups []JsonBloodGroup `json:"BloodGroups"`
}

// --- Seeding Functions ---

//  Seed Cities & Regions
func SeedMoroccanCities(db *gorm.DB) {
	var count int64
	db.Model(&models.Region{}).Count(&count)
	if count > 0 {
		fmt.Println("🇲🇦 Moroccan Data already exists. Skipping.")
		return
	}

	fmt.Println("⏳ Seeding Regions & Cities...")

	// 1. Regions
	dataRegion, err := os.ReadFile("assets/sql-moroccan-cities/json/region.json")
	if err != nil {
		fmt.Println("❌ Error reading region.json:", err)
		return
	}

	var jsonRegions []JsonRegion
	json.Unmarshal(dataRegion, &jsonRegions)

	regionMap := make(map[string]uint)

	for _, r := range jsonRegions {
		region := models.Region{Region: r.Region}
		if err := db.Create(&region).Error; err == nil {
			regionMap[r.ID] = region.ID // 
		}
	}
	fmt.Println("✅ Regions Imported.")

	//  Cities
	dataCity, err := os.ReadFile("assets/sql-moroccan-cities/json/ville.json")
	if err != nil {
		fmt.Println("❌ Error reading ville.json:", err)
		return
	}

	var jsonCities []JsonCity
	json.Unmarshal(dataCity, &jsonCities)

	var citiesToInsert []models.City
	for _, c := range jsonCities {
		realRegionID, exists := regionMap[c.RegionID]
		if exists {
			citiesToInsert = append(citiesToInsert, models.City{
				Ville:    c.Ville,
				RegionID: realRegionID,
			})
		}
	}

	if len(citiesToInsert) > 0 {
		db.CreateInBatches(citiesToInsert, 100)
	}
	fmt.Println("✅ Cities Imported Successfully.")
}

//  Seed Banks
func SeedBanks(db *gorm.DB) {
	var count int64
	db.Model(&models.Bank{}).Count(&count)
	if count > 0 {
		fmt.Println("🏥 Banks already seeded.")
		return
	}

	fmt.Println("⏳ Seeding Hospitals Data...")

	data, err := os.ReadFile("assets/banks.json")
	if err != nil {
		fmt.Println("❌ Error reading banks.json:", err)
		return
	}

	var jsonBanks []JsonBank
	if err := json.Unmarshal(data, &jsonBanks); err != nil {
		fmt.Println("❌ Error parsing JSON:", err)
		return
	}

	for _, b := range jsonBanks {
		// Prepare Stock
		var stock []models.BloodGroup

		stockMap := make(map[string]int)
		for _, bg := range b.BloodGroups {
			stockMap[bg.Group] = bg.Value
		}

		allGroups := []string{"A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"}
		for _, g := range allGroups {
			quantity := 0
			if val, ok := stockMap[g]; ok {
				quantity = val
			}
			stock = append(stock, models.BloodGroup{
				Group:    g,
				Quantity: quantity,
			})
		}

		// Create Bank Object
		bank := models.Bank{
			Name:        b.Name,
			Email:       b.Email,
			PhoneNumber: b.PhoneNumber,
			Website:     b.Website,
			// GORM غايتكلف بإنشاء Address و BloodGroups أوتوماتيكيا
			Address: models.Address{
				Area:       b.Address.Area,
				City:       b.Address.City,
				State:      b.Address.State,
				Country:    "Morocco",
				PostalCode: b.Address.PostalCode,
			},
			BloodGroups: stock,
		}

		if err := db.Create(&bank).Error; err != nil {
			fmt.Printf("⚠️ Failed to seed bank %s: %v\n", b.Name, err)
		}
	}
	fmt.Println("✅ Hospitals Seeded Successfully.")
}
